package xfn

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PackagedFunctionRunner struct {
	client       client.Reader
	creds        credentials.TransportCredentials
	interceptors []InterceptorCreator
}

type InterceptorCreator interface {
	CreateInterceptor(name, okg string) grpc.UnaryClientInterceptor
}
