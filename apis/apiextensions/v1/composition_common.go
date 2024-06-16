package v1

// A CompositionMode determines what mode of Composition is used.
// shouldnt it be like an enum with pre-determined values?
type CompositionMode string

// TypeReference is used to refer to a type for declaring compatibility.
type TypeReference struct {
	// APIVersion of the type.
	APIVersion string `json:"apiVersion"`

	// Kind of the type.
	Kind string `json:"kind"`
}
