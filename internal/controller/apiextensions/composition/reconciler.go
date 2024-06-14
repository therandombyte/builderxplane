// Creates composition revisions

package composition

import (
	"context"
	"strings"
	"time"

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

const (
	timeout = 2 * time.Minute
)

// Error strings.
const (
	errGet             = "cannot get Composition"
	errListRevs        = "cannot list CompositionRevisions"
	errCreateRev       = "cannot create CompositionRevision"
	errOwnRev          = "cannot own CompositionRevision"
	errUpdateRevStatus = "cannot update CompositionRevision status"
	errUpdateRevSpec   = "cannot update CompositionRevision spec"
)

// Reconciliation logic
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := r.log.WithValues("request", req)
	log.Info("Reconciling via Debug")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	comp := &v1.Composition{}
	if err := r.client.Get(ctx, req.NamespacedName, comp); err != nil {
		log.Info(errGet, "error", err)
		return reconcile.Result{}, nil
	}

	var latestRev int64

	if err := r.client.Create(ctx, NewCompositionRevision(comp, latestRev+1)); err != nil {
		log.Info(errCreateRev, "error", err)
		return reconcile.Result{}, nil
	}
	log.Info("Created new revision", "revision", latestRev+1)
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
