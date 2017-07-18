package gohost

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

// ServeGRPCWithTLS starts a gRPC endpoint for the given server with TLS enabled.
func ServeGRPCWithTLS(server Server, grpcAddr string, opts []grpc.ServerOption, certFile string, keyFile string) error {
	// create TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to generate TLS credentials: %v", err)
	}

	// add TLS credentials to options
	opts = append(opts, grpc.Creds(creds))

	// start server
	return ServeGRPC(server, grpcAddr, opts)
}

// ServeHTTP starts an HTTP endpoint for a given server. This is a gateway pointing to a gRPC endpoint.
func ServeHTTP(server Server, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption) error {
	// start server
	return serveHTTP(server, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler http.Handler) error {
		return http.ListenAndServe(addr, handler)
	})
}

// ServeHTTPWithTLS starts an HTTP endpoint for a given server with TLS enabled. This is a gateway pointing to a gRPC endpoint.
func ServeHTTPWithTLS(server Server, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption, certFile string, keyFile string, insecureSkipVerify bool) error {
	// create TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	})

	// add TLS credentials to options
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// start server
	return serveHTTP(server, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler http.Handler) error {
		return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	})
}

// serveHTTP starts an HTTP endpoint for a given server. This is a gateway pointing to a gRPC endpoint.
func serveHTTP(server Server, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption, listenAndServe func(addr string, handler http.Handler) error) error {
	// create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// register server
	var handler http.Handler
	handler = runtime.NewServeMux()
	err := server.RegisterHandler(ctx, handler, grpcAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register HTTP endpoint: %v", err)
	}

	// enable CORS if requested
	if enableCORS {
		handler = cors.AllowAll().Handler(handler)
	}

	// start server
	err = listenAndServe(httpAddr, handler)
	if err != nil {
		return fmt.Errorf("failed to serve HTTP endpoint: %v", err)
	}

	return nil
}
