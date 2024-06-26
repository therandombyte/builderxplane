package ratelimiter

import (
	"context"
	"sync"
	"time"

	"github.com/therandombyte/builderxplane/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"
)

// a wrapper reconciler to rate limit
type Reconciler struct {
	name  string
	inner reconcile.Reconciler
	limit ratelimiter.RateLimiter

	limited  map[string]struct{}
	limitedL sync.RWMutex
}

func NewReconciler(name string, r reconcile.Reconciler, l ratelimiter.RateLimiter) *Reconciler {
	return &Reconciler{name: name, inner: r, limit: l, limited: make(map[string]struct{})}
}

// Reconcile the supplied request subject to rate limiting.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	item := r.name + req.String()
	if d := r.when(req); d > 0 {
		return reconcile.Result{RequeueAfter: d}, nil
	}
	r.limit.Forget(item)
	return r.inner.Reconcile(ctx, req)
}

// when adapts the upstream rate limiter's 'When' method such that rate limited
// requests can call it again when they return and will be allowed to proceed
// immediately without being subject to further rate limiting. It is optimised
// for handling requests that have not been and will not be rate limited without
// blocking.
func (r *Reconciler) when(req reconcile.Request) time.Duration {
	item := r.name + req.String()

	r.limitedL.RLock()
	_, limited := r.limited[item]
	r.limitedL.RUnlock()

	// If we already rate limited this request we trust that it complied and
	// let it pass immediately.
	if limited {
		r.limitedL.Lock()
		delete(r.limited, item)
		r.limitedL.Unlock()
		return 0
	}

	d := r.limit.When(item)

	// Record that this request was rate limited so that we can let it
	// through immediately when it requeues after the supplied duration.
	if d != 0 {
		r.limitedL.Lock()
		r.limited[item] = struct{}{}
		r.limitedL.Unlock()
	}

	return d
}
