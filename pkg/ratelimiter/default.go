package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"
)

// talking to client-go, probably thats why the runtime is separated
// overriding default k8s controller manager reconcile rate from 30 to 10
func LimitRESTConfig(cfg *rest.Config, rps int) *rest.Config {
	out := rest.CopyConfig(cfg)
	out.QPS = float32(rps * 5)
	out.Burst = rps * 10
	return out
}

// setting the max requeue size per second
func NewGlobal(rps int) *workqueue.BucketRateLimiter {
	return &workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(rps), rps*10)}
}

func NewController() ratelimiter.RateLimiter {
	return workqueue.NewItemExponentialFailureRateLimiter(1*time.Second, 60*time.Second)
}
