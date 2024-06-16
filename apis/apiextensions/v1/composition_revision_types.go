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
// +kubebuilder:storageversion
// +genclient
// +genclient:nonNamespaced

// A CompositionRevision represents a revision of a Composition. Crossplane
// creates new revisions when there are changes to the Composition.
//
// Crossplane creates and manages CompositionRevisions. Don't directly edit
// CompositionRevisions.
// +kubebuilder:printcolumn:name="REVISION",type="string",JSONPath=".spec.revision"
// +kubebuilder:printcolumn:name="XR-KIND",type="string",JSONPath=".spec.compositeTypeRef.kind"
// +kubebuilder:printcolumn:name="XR-APIVERSION",type="string",JSONPath=".spec.compositeTypeRef.apiVersion"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories=crossplane,shortName=comprev
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
