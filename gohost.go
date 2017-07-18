package gohost

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// ServeGRPC starts a gRPC endpoint for the given server.
func ServeGRPC(server Server, grpcAddr string, opts []grpc.ServerOption) error {
	// start listening
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}

	// register server
	grpcServer := grpc.NewServer(opts...)
	server.RegisterServer(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("Failed to serve: %v", err)
	}

	return nil
}
