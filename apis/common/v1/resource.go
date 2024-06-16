package v1

import "k8s.io/apimachinery/pkg/types"

type TypedReference struct {
	APIVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	Name       string    `json:"name"`
	UID        types.UID `json:"uid,omitempty"`
}
