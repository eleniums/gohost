package gohost

import (
	"google.golang.org/grpc"
)

// GRPCService is the interface for a hosted service with a gRPC endpoint.
type GRPCService interface {
	// RegisterServer registers this server to be a gRPC endpoint.
	RegisterServer(grpc *grpc.Server)
}
