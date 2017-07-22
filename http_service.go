package gohost

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// HTTPService is the interface for a hosted service with a HTTP endpoint.
type HTTPService interface {
	// RegisterHandler registers the HTTP handler to use with a service.
	RegisterHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}
