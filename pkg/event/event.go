// Records k8s events
package event

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type Recorder interface {
	Event(obj runtime.Object, e Event)
	WithAnnotations(keysAndValues ...string) Recorder
}

// type pf event
type Type string

type Reason string

// xplane resource event
type Event struct {
	Type        Type
	Reason      Reason
	Message     string
	Annotations map[string]string
}

// An APIRecorder records Kubernetes events to an API server.
type APIRecorder struct {
	kube        record.EventRecorder
	annotations map[string]string
}

// NewAPIRecorder returns an APIRecorder that records Kubernetes events to an
// APIServer using the supplied EventRecorder.
func NewAPIRecorder(r record.EventRecorder) *APIRecorder {
	return &APIRecorder{kube: r, annotations: map[string]string{}}
}

// Event records the supplied event.
func (r *APIRecorder) Event(obj runtime.Object, e Event) {
	r.kube.AnnotatedEventf(obj, r.annotations, string(e.Type), string(e.Reason), e.Message)
}

// WithAnnotations returns a new *APIRecorder that includes the supplied
// annotations with all recorded events.
func (r *APIRecorder) WithAnnotations(keysAndValues ...string) Recorder {
	ar := NewAPIRecorder(r.kube)
	for k, v := range r.annotations {
		ar.annotations[k] = v
	}
	sliceMap(keysAndValues, ar.annotations)
	return ar
}

func sliceMap(from []string, to map[string]string) {
	for i := 0; i+1 < len(from); i += 2 {
		k, v := from[i], from[i+1]
		to[k] = v
	}
}

// A NopRecorder does nothing.
type NopRecorder struct{}

// NewNopRecorder returns a Recorder that does nothing.
func NewNopRecorder() *NopRecorder {
	return &NopRecorder{}
}

// Event does nothing.
func (r *NopRecorder) Event(_ runtime.Object, _ Event) {}

// WithAnnotations does nothing.
func (r *NopRecorder) WithAnnotations(_ ...string) Recorder { return r }
