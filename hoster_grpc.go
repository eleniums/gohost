package gohost

import (
	"errors"
	"fmt"
	"net"

	rpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// serveGRPC will start the gRPC endpoint.
func (h *Hoster) serveGRPC() error {
	// configure server options
	serverOpts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(h.MaxSendMsgSize),
		grpc.MaxRecvMsgSize(h.MaxRecvMsgSize),
	}

	// add interceptors
	if len(h.UnaryInterceptors) > 0 {
		unaryInterceptorChain := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(h.UnaryInterceptors...))
		serverOpts = append(serverOpts, unaryInterceptorChain)
	}
	if len(h.StreamInterceptors) > 0 {
		streamInterceptorChain := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(h.StreamInterceptors...))
		serverOpts = append(serverOpts, streamInterceptorChain)
	}

	// start the gRPC endpoint
	if h.IsTLSEnabled() {
		return gogrpc.ServeGRPCWithTLS(h.Service, h.GRPCAddr, serverOpts, h.CertFile, h.KeyFile)
	}

	return gogrpc.ServeGRPC(h.Service, h.GRPCAddr, serverOpts)
}

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
