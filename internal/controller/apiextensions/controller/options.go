package controller

import (
	"github.com/therandombyte/builderxplane/internal/engine"
	"github.com/therandombyte/builderxplane/internal/xfn"
	"github.com/therandombyte/builderxplane/pkg/controller"
)

type Options struct {
	controller.Options

	ControllerEngine *engine.ControllerEngine
	FunctionRunner   *xfn.PackagedFunctionRunner
}
