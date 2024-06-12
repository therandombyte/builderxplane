// package to configure controller options

package controller

import (
	"crypto/tls"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/statemetrics"
	"github.com/therandombyte/builderxplane/pkg/feature"
	"github.com/therandombyte/builderxplane/pkg/logging"
	"github.com/therandombyte/builderxplane/pkg/ratelimiter"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

type Options struct {
	Logger                  logging.Logger
	GlobalRateLimiter       workqueue.RateLimiter // reconcilation rate
	PollInterval            time.Duration         // if contorller has work to do
	MaxConcurrentReconciles int
	Features                *feature.Flags // a map of features (enabled or not)
	ESSOption               *ESSOptions
	MetricOptions           *MetricOptions
}

type ESSOptions struct {
	TLSConfig     *tls.Config
	TLSSecretname *string
}

// Metrics not impemented locally
type MetricOptions struct {
	PollStateMetricInterval time.Duration
	MRMetrics               managed.MetricRecorder
	MRStateMetrics          *statemetrics.MRStateMetrics
}

// extract options for controller-runtime
func (o Options) ForControllerRuntime() controller.Options {
	recoverPanic := true
	return controller.Options{
		MaxConcurrentReconciles: o.MaxConcurrentReconciles,
		RateLimiter:             ratelimiter.NewController(),
		RecoverPanic:            &recoverPanic,
	}
}
