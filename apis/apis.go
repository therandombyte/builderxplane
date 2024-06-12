// contains k8s API groups

package apis

import (
	"github.com/therandombyte/builderxplane/apis/apiextensions"
	"k8s.io/apimachinery/pkg/runtime"
)

// Register the types with the scheme so that components can map to GVK and back
func init() {
	AddToSchemes = append(AddToSchemes, apiextensions.AddToScheme)
}

// to add all resources defined in xplane to a Scheme
var AddToSchemes runtime.SchemeBuilder

func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
