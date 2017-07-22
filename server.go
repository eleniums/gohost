package gohost

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server is the interface for a hosted service.
type Server interface {
	// RegisterServer registers this server to be a gRPC endpoint.
	RegisterServer(grpc *grpc.Server)

	// RegisterHandler registers this server to be an HTTP endpoint.
	RegisterHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}
