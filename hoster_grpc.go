package gohost

import (
	"errors"
	"fmt"
	"net"

	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ServeGRPC starts a gRPC endpoint for the given service.
func ServeGRPC(service GRPCService, grpcAddr string, opts []rpc.ServerOption) error {
	// validate parameters
	if service == nil {
		return errors.New("service cannot be nil")
	}
	if grpcAddr == "" {
		return errors.New("grpc address cannot be empty")
	}

	// start listening
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// register server
	grpcServer := rpc.NewServer(opts...)
	service.RegisterServer(grpcServer)

	// start server
	return grpcServer.Serve(lis)
}

// ServeGRPCWithTLS starts a gRPC endpoint for the given service with TLS enabled.
func ServeGRPCWithTLS(service GRPCService, grpcAddr string, opts []rpc.ServerOption, certFile string, keyFile string) error {
	// validate parameters
	if certFile == "" {
		return errors.New("cert file cannot be empty")
	}
	if keyFile == "" {
		return errors.New("key file cannot be empty")
	}

	// create TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS credentials: %v", err)
	}

	// add TLS credentials to options
	opts = append(opts, rpc.Creds(creds))

	// start server
	return ServeGRPC(service, grpcAddr, opts)
}
