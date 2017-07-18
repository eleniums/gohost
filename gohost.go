package gohost

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

// ServeGRPC starts a gRPC endpoint for the given server.
func ServeGRPC(server Server, grpcAddr string, opts []grpc.ServerOption) error {
	// start listening
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// register server
	grpcServer := grpc.NewServer(opts...)
	server.RegisterServer(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC endpoint: %v", err)
	}

	return nil
}

// ServeHTTP starts an HTTP endpoint for a given server. This is a gateway pointing to a gRPC endpoint.
func ServeHTTP(server Server, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption) error {
	// create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// register server
	mux := runtime.NewServeMux()
	err := server.RegisterHandler(ctx, mux, grpcAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register HTTP endpoint: %v", err)
	}

	// enable CORS if requested
	if enableCORS {
		mux := cors.AllowAll().Handler(mux)
	}

	// start server
	err = http.ListenAndServe(httpAddr, mux)
	if err != nil {
		return fmt.Errorf("failed to serve HTTP endpoint: %v", err)
	}

	return nil
}
