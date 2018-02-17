package gohost

import (
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type grpcServer func(s *grpc.Server)

// serveGRPC will start the gRPC endpoint.
func (h *Hoster) serveGRPC() error {
	// validate parameters
	if len(h.grpcEndpoints) == 0 {
		return errors.New("no grpc servers added")
	}
	if h.GRPCAddr == "" {
		return errors.New("grpc address cannot be empty")
	}

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
		return ServeGRPCWithTLS(h.grpcEndpoints, h.GRPCAddr, serverOpts, h.CertFile, h.KeyFile)
	}

	return ServeGRPC(h.grpcEndpoints, h.GRPCAddr, serverOpts)
}

// TODO: move this code up into the above method
// ServeGRPC starts a gRPC endpoint for the given service.
func ServeGRPC(servers []grpcServer, addr string, opts []grpc.ServerOption) error {
	// start listening
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// register servers
	grpcServer := grpc.NewServer(opts...)
	for i := range servers {
		servers[i](grpcServer)
	}

	// start servers
	return grpcServer.Serve(lis)
}

// ServeGRPCWithTLS starts a gRPC endpoint for the given service with TLS enabled.
func ServeGRPCWithTLS(servers []grpcServer, addr string, opts []grpc.ServerOption, certFile string, keyFile string) error {
	// create TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS credentials: %v", err)
	}

	// add TLS credentials to options
	opts = append(opts, grpc.Creds(creds))

	// start server
	return ServeGRPC(servers, addr, opts)
}
