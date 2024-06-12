package engine

import (
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// a controller watch source that can be stopped
type StoppableSource struct {
	infs cache.Informers

	Type       client.Object
	handler    handler.EventHandler
	predicates []predicate.Predicate
}
