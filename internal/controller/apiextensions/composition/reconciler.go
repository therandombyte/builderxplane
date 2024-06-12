// Creates composition revisions

package composition

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/therandombyte/builderxplane/apis/apiextensions/v1"
	"github.com/therandombyte/builderxplane/internal/controller/apiextensions/controller"
	"github.com/therandombyte/builderxplane/pkg/event"
	"github.com/therandombyte/builderxplane/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrl "sigs.k8s.io/controller-runtime"
)

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := "/revisions" + strings.ToLower(v1.CompositionGroupKind)

	r := NewReconciler(mgr,
		WithLogger(o.Logger.WithValues("controller", name)),
		WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1.Composition{}).
		Owns(&v1.CompositionRevision{}).
		WithOptions(o.ForControllerRuntime()).
		Complete(r)
}

// reconciliation by creating new CompositionRevisions for each revision of spec
type Reconciler struct {
	client client.Client
	log    logging.Logger
	record event.Recorder
}

// ReconcilerOption is used to configure the Reconciler.
type ReconcilerOption func(*Reconciler)

func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", req)
	log.Debug("Reconciling via Debug")
	fmt.Println("Reconciling via Println")

	return reconcile.Result{}, nil
}

func NewReconciler(mgr manager.Manager, opts ...ReconcilerOption) *Reconciler {
	r := &Reconciler{
		client: mgr.GetClient(),
		log:    logging.NewNopLogger(),
		record: event.NewNopRecorder(),
	}

	for _, f := range opts {
		f(r)
	}
	return r
}

// WithLogger specifies how the Reconciler should log messages.
func WithLogger(log logging.Logger) ReconcilerOption {
	return func(r *Reconciler) {
		r.log = log
	}
}

// WithRecorder specifies how the Reconciler should record Kubernetes events.
func WithRecorder(er event.Recorder) ReconcilerOption {
	return func(r *Reconciler) {
		r.record = er
	}
}
