package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CompositionRevisionSpec struct {
	CompositeTypeRef TypeReference `json:"compositeTypeRef,omitempty"`
}

type CompositionRevisionStatus struct {
	//xpv1.ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type CompositionRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompositionRevisionSpec   `json:"spec,omitempty"`
	Status CompositionRevisionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CompositionRevisionList contains a list of CompositionRevisions.
type CompositionRevisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CompositionRevision `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CompositionRevision{}, &CompositionRevisionList{})
}
