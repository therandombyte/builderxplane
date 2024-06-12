package reconcile

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/types"
)

// contains info (name and namespace) to reconcile a k8s object
type Request struct {
	types.NamespacedName
}

// result of a reconcile operation
type Result struct {
	// tell the controller to requeue the reconcile key
	Requeue      bool
	RequeueAfter time.Duration
}

type Reconciler interface {
	Reconcile(context.Context, Request) (Result, error)
}

// Func is a func that implements reconcile Interface
type Func func(context.Context, Request) (Result, error)

var _ Reconciler = Func(nil)

func (r Func) Reconcile(ctx context.Context, o Request) (Result, error) { return r(ctx, o) }
