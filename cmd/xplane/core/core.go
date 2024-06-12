// Implementing core controller manager
package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/therandombyte/builderxplane/internal/controller/apiextensions"
	apiextensionscontroller "github.com/therandombyte/builderxplane/internal/controller/apiextensions/controller"
	"github.com/therandombyte/builderxplane/internal/engine"
	"github.com/therandombyte/builderxplane/internal/xfn"
	"github.com/therandombyte/builderxplane/pkg/controller"
	"github.com/therandombyte/builderxplane/pkg/feature"
	"github.com/therandombyte/builderxplane/pkg/logging"
	"github.com/therandombyte/builderxplane/pkg/ratelimiter"
	"github.com/therandombyte/builderxplane/pkg/resource/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// this is the export
type Command struct {
	Start startCommand `cmd:"" help:"Start xplane controllers"`
}

// keep this internal, not to export
type startCommand struct {
	Profile string `help:"Serve runtime profiling via HTTP at /debug/pprof." placeholder:"host:port"`

	Namespace         string        `default:"xplane-system"`
	SyncInterval      time.Duration `default:"1h" help:"How often all resources will be double-checked for drift from desired state"`
	TLSServerCertsDir string        `env:"TLS_SERVER_CERTS_DIR"  help:"Folder path that stored TLS certs of xplane"`
	LeaderElection    bool          `default:"false" env:"LEADER_ELECTION"`
	MaxReconcileRate  int           `default:"100" help:"max rpm of drift check"`
	PollInterval      time.Duration `default:"1m" help:"drift check rate for individual resources"`
}

// root will be executed at the last
func (c *Command) Run() error {
	fmt.Println("In Run of Root start command")
	return nil
}

// logger is from crossplane runtime, shall we remove the dependency
func (c *startCommand) Run(s *runtime.Scheme, log logging.Logger) error {
	fmt.Println("In Run of Core start command")

	// 1. GetConfig() config holds common attributes that can be passed to k8s client on init
	// like QPS, ratelimiter, warninghandler etc
	cfg, err := ctrl.GetConfig()
	if err != nil {
		fmt.Printf("Error getting config: %s", err.Error())
		return nil
	}

	// 2. NewBroadcaster()
	eb := record.NewBroadcaster()

	// 3. NewManager(): adding a ratelimiter to customize the restconfig, creates a copy of config
	// and changes QPS and Burst
	mgr, err := ctrl.NewManager(ratelimiter.LimitRESTConfig(cfg, 10), ctrl.Options{
		Scheme: s,
		Cache: cache.Options{
			SyncPeriod: &c.SyncInterval, // default to 1hr
		},
		WebhookServer: webhook.NewServer(webhook.Options{ // go http server
			CertDir: c.TLSServerCertsDir,
			TLSOpts: []func(*tls.Config){
				func(t *tls.Config) {
					t.MinVersion = tls.VersionTLS13
				},
			},
		}),
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor:   []client.Object{},
				Unstructured: false,
			},
		},
		EventBroadcaster:              eb,
		LeaderElection:                c.LeaderElection,
		LeaderElectionID:              "xplane-leader-election-core",
		LeaderElectionResourceLock:    resourcelock.LeasesResourceLock,
		LeaderElectionReleaseOnCancel: true,
		LeaseDuration:                 func() *time.Duration { d := 60 * time.Second; return &d }(),
		RenewDeadline:                 func() *time.Duration { d := 50 * time.Second; return &d }(),

		PprofBindAddress:       c.Profile,
		HealthProbeBindAddress: ":8081",
	})

	if err != nil {
		fmt.Printf("Cannot create manager %s", err.Error())
	}

	// no use yet
	eb.StartLogging(func(format string, args ...interface{}) {
		log.Debug(fmt.Sprintf(format, args...))
	})
	defer eb.Shutdown()

	// no use yet
	o := controller.Options{
		Logger:                  log,
		MaxConcurrentReconciles: c.MaxReconcileRate,
		PollInterval:            c.PollInterval,
		GlobalRateLimiter:       ratelimiter.NewGlobal(c.MaxReconcileRate),
		Features:                &feature.Flags{},
	}

	// 4. Cache for claim and XR controller
	ca, err := cache.New(mgr.GetConfig(), cache.Options{
		HTTPClient: mgr.GetHTTPClient(),
		Scheme:     mgr.GetScheme(),
		Mapper:     mgr.GetRESTMapper(),
		SyncPeriod: &c.SyncInterval,
	})
	if err != nil {
		fmt.Printf("cannot create cache for API extensions controller %s", err.Error())
	}

	// 5. create context for cache start
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		// dont start the cache untill manager is elected
		<-mgr.Elected()
		if err := ca.Start(ctx); err != nil {
			log.Info("API extensions cache returned error", "error", err)
			fmt.Println("API extensions cache returned error", "error", err)
		}
		fmt.Println("API extensions cache stopped")
	}()

	// 6. Create Client to create the controller engine
	cl, err := client.New(mgr.GetConfig(), client.Options{
		HTTPClient: mgr.GetHTTPClient(),
		Scheme:     mgr.GetScheme(),
		Mapper:     mgr.GetRESTMapper(),
		Cache: &client.CacheOptions{
			Reader:       ca,
			Unstructured: true,
		},
	})

	// 7. Create engine
	ce := engine.New(mgr,
		engine.TrackInformers(ca, mgr.GetScheme()),
		unstructured.NewClient(cl),
		engine.WithLogger(log),
	)

	// 8. Creating API Controller Options
	var functionRunner *xfn.PackagedFunctionRunner
	ao := apiextensionscontroller.Options{
		Options:          o,
		ControllerEngine: ce,
		FunctionRunner:   functionRunner,
	}

	if err := apiextensions.Setup(mgr, ao); err != nil {
		fmt.Printf("Cannot create API extension controllers: %s", err.Error())
	}

	mgr.Start(ctrl.SetupSignalHandler())

	fmt.Println("Controller Options", o)
	fmt.Println("Creating Engine =====", ce)
	//cfg.WarningHandler =
	return nil
}
