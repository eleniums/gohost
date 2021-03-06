package gohost

import (
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// serveGRPC will start the gRPC endpoint.
func (h *Hoster) serveGRPC() error {
	// validate parameters
	if h.GRPCAddr == "" {
		return errors.New("grpc address cannot be empty")
	}

	// configure server options
	opts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(h.MaxSendMsgSize),
		grpc.MaxRecvMsgSize(h.MaxRecvMsgSize),
	}

	// add interceptors
	if len(h.UnaryInterceptors) > 0 {
		unaryInterceptorChain := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(h.UnaryInterceptors...))
		opts = append(opts, unaryInterceptorChain)
	}
	if len(h.StreamInterceptors) > 0 {
		streamInterceptorChain := grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(h.StreamInterceptors...))
		opts = append(opts, streamInterceptorChain)
	}

	// add TLS credentials to options if necessary
	if h.isTLSEnabled() {
		creds, err := credentials.NewServerTLSFromFile(h.CertFile, h.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS credentials: %v", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	// start listening
	lis, err := net.Listen("tcp", h.GRPCAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// register servers
	server := grpc.NewServer(opts...)
	for i := range h.grpcServers {
		h.grpcServers[i](server)
	}

	// start the gRPC endpoint
	return server.Serve(lis)
}

// isTLSEnabled will return true if TLS properties are set and ready to use.
func (h *Hoster) isTLSEnabled() bool {
	return h.CertFile != "" && h.KeyFile != ""
}
