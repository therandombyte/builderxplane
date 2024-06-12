package engine

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type InformerTrackingCache struct {
	cache.Cache
	scheme *runtime.Scheme
	mx     sync.RWMutex
	active map[schema.GroupVersionKind]bool
}

// TrackInformers wraps the supplied cache, adding a method to query which
// informers are active.
func TrackInformers(c cache.Cache, s *runtime.Scheme) *InformerTrackingCache {
	return &InformerTrackingCache{
		Cache:  c,
		scheme: s,
		active: make(map[schema.GroupVersionKind]bool),
	}
}

// returns GVK of active Informers
func (c *InformerTrackingCache) ActiveInformers() []schema.GroupVersionKind {
	c.mx.RLock()
	defer c.mx.RUnlock()

	out := make([]schema.GroupVersionKind, 0, len(c.active))
	for gvk := range c.active {
		out = append(out, gvk)
	}
	return out
}
