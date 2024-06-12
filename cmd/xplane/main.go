package main

import (
	"fmt"
	"io"

	"github.com/alecthomas/kong"
	"github.com/therandombyte/builderxplane/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/therandombyte/builderxplane/apis"
	apiextensionsv1 "github.com/therandombyte/builderxplane/apis/apiextensions/v1"
	"github.com/therandombyte/builderxplane/cmd/xplane/core"
	admv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
)

// -------------- Picked up from Kubebuilder -------------
var (
	scheme = runtime.NewScheme()
	xlog   = ctrl.Log.WithName("xplane") // xplane uses logr and klog
)

func init() {
	fmt.Println("=========== IN MAIN INIT ============")
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

// --------------------------------------------------------

// defining custom types
type (
	debugFlag   bool
	versionFlag bool
)

type cli struct {
	Debug   debugFlag   `help:"Print verbose logging statements." short:"d"`
	Version versionFlag `help:"Print version and quit" short:"v"`

	Core core.Command `cmd:"" default:"withargs" help:"Start core xplane contollers"`
}

func main() {
	z1 := zap.New().WithName("xplane") // for JSON type logging, in k8s implmentaion, it includes GVK
	logging.SetFilteredKlogLogger(z1)  // setting log backend
	z1.Info("Main of xplane")

	ctrl.SetLogger(zap.New(zap.WriteTo((io.Discard)))) // prevent excessive logging

	// all custom Go objects get added to scheme and later gets converted to GVK
	s := runtime.NewScheme()

	/*
		Break leaf commands out into separate structs.
		Attach a Run(...) error method to all leaf commands.
		Call kong.Kong.Parse() to obtain a kong.Context.
		Call kong.Context.Run(bindings...) to call the selected parsed command.
	*/
	ctx := kong.Parse(&cli{},
		kong.Name("xplane"),
		kong.Description("An open source multicloud control plane."),
		kong.BindTo(logging.NewLogrLogger(z1), (*logging.Logger)(nil)),
		kong.UsageOnError(),
	)

	ctx.FatalIfErrorf(corev1.AddToScheme(s), "cannot add core v1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(appsv1.AddToScheme(s), "cannot add apps v1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(rbacv1.AddToScheme(s), "cannot add rbac v1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(coordinationv1.AddToScheme(s), "cannot add coordination v1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(extv1.AddToScheme(s), "cannot add apiextensions v1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(extv1beta1.AddToScheme(s), "cannot add apiextensions v1beta1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(admv1.AddToScheme(s), "cannot add admissionregistration v1 Kubernetes API types to scheme")
	ctx.FatalIfErrorf(apis.AddToScheme(s), "cannot add Crossplane API types to scheme")
	ctx.FatalIfErrorf(ctx.Run(s))
}





/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
/*
package main

import (
	"crypto/tls"
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	apiextensionsv1 "github.com/therandombyte/builderxplane/api/v1"
	"github.com/therandombyte/builderxplane/internal/controller"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var secureMetrics bool
	var enableHTTP2 bool
	flag.StringVar(&metricsAddr, "metrics-bind-address", "0", "The address the metric endpoint binds to. "+
		"Use the port :8080. If not set, it will be 0 in order to disable the metrics server")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", false,
		"If set the metrics endpoint is served securely")
	flag.BoolVar(&enableHTTP2, "enable-http2", false,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// if the enable-http2 flag is false (the default), http/2 should be disabled
	// due to its vulnerabilities. More specifically, disabling http/2 will
	// prevent from being vulnerable to the HTTP/2 Stream Cancellation and
	// Rapid Reset CVEs. For more information see:
	// - https://github.com/advisories/GHSA-qppj-fm5r-hxr3
	// - https://github.com/advisories/GHSA-4374-p667-p6c8
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	tlsOpts := []func(*tls.Config){}
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			TLSOpts:       tlsOpts,
		},
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "940fa6d6.xplane.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controller.CompositionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Composition")
		os.Exit(1)
	}
	if err = (&controller.CompositionRevisionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CompositionRevision")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

*/
