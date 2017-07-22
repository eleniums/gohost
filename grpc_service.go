package gohost

import (
	"google.golang.org/grpc"
)

// GRPCService is the interface for a hosted service with a gRPC endpoint.
type GRPCService interface {
	// RegisterServer registers the gRPC server to use with a service.
	RegisterServer(grpc *grpc.Server)
}
