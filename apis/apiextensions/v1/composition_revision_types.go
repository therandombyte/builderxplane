package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// LabelCompositionName is the name of the Composition used to create
	// this CompositionRevision.
	LabelCompositionName = "xpane.io/composition-name"

	// LabelCompositionHash is a hash of the Composition label, annotation
	// and spec used to create this CompositionRevision. Used to identify
	// identical revisions.
	LabelCompositionHash = "xplane.io/composition-hash"
)

type CompositionRevisionSpec struct {
	CompositeTypeRef TypeReference    `json:"compositeTypeRef,omitempty"`
	Mode             *CompositionMode `json:"mode,omitempty"`
	Revision         int64            `json:"revision"`
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
