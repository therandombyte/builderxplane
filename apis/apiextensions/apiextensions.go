// contains k8s API groups for extension type of xplane

package apiextensions

import (
	v1 "github.com/therandombyte/builderxplane/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	AddToSchemes = append(AddToSchemes, v1.AddToScheme)
}

var AddToSchemes runtime.SchemeBuilder

func AddToScheme(s *runtime.Scheme) error {
	return AddToSchemes.AddToScheme(s)
}
