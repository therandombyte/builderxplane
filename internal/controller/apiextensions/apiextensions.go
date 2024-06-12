// Implements xplane composition controllers

package apiextensions

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/therandombyte/builderxplane/internal/controller/apiextensions/composition"
	"github.com/therandombyte/builderxplane/internal/controller/apiextensions/controller"
)

// setup API extensions controller
func Setup(mgr ctrl.Manager, o controller.Options) error {
	return composition.Setup(mgr, o)
}
