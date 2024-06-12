// manages the lifecycle of controllers

package engine

import (
	"context"
	"sync"

	"github.com/therandombyte/builderxplane/pkg/logging"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	kcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// has its own client and cache, does not use the manager's
type ControllerEngine struct {
	mgr    manager.Manager
	infs   TrackingInformers
	client client.Client

	log         logging.Logger
	mx          sync.RWMutex
	controllers map[string]*controller
}

type TrackingInformers interface {
	cache.Informers
	ActiveInformers() []schema.GroupVersionKind
}

type controller struct {
	ctrl    kcontroller.Controller
	cancel  context.CancelFunc
	mx      sync.RWMutex
	sources map[WatchID]*StoppableSource
}

type WatchType string

// A WatchID uniquely identifies a watch.
type WatchID struct {
	Type WatchType
	GVK  schema.GroupVersionKind
}

type ControllerEngineOption func(*ControllerEngine)

func New(mgr manager.Manager, infs TrackingInformers, c client.Client, o ...ControllerEngineOption) *ControllerEngine {
	e := &ControllerEngine{
		mgr:         mgr,
		infs:        infs,
		client:      c,
		log:         logging.NewNopLogger(),
		controllers: make(map[string]*controller),
	}

	for _, fn := range o {
		fn(e)
	}

	return e
}

func WithLogger(l logging.Logger) ControllerEngineOption {
	return func(e *ControllerEngine) {
		e.log = l
	}
}
