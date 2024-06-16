package v1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

const (
	Group   = "apiextensions.xplane.io"
	Version = "v1"
)

var (
	// SchemGroupVersion is the group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

	// SchemaBuilder is used to add go types to the GroupVersionKind scheme
	SchemaBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	// AddToScheme adds all registered types to scheme
	//AddToScheme = SchemaBuilder.AddToScheme
)

// Composition type metadata
var (
	CompositionKind             = reflect.TypeOf(Composition{}).Name()
	CompositionGroupKind        = schema.GroupKind{Group: Group, Kind: CompositionKind}.String()
	CompositionKindAPIVersion   = CompositionKind + "." + SchemeGroupVersion.String()
	CompositionGroupVersionKind = SchemeGroupVersion.WithKind(CompositionKind)
)
