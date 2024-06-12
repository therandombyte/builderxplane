package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CompositionSpec struct {
	CompositeTypeRef TypeReference `json:"compositeTypeRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type Composition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CompositionSpec `json:"spec,omitempty"`
}

// CompositionStatus defines the observed state of Composition
type CompositionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

type CompositionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Composition `json:"items"`
}

func init() {
	fmt.Println("=========== IN COMPOSITION INIT ============")
	SchemeBuilder.Register(&Composition{}, &CompositionList{})
}
