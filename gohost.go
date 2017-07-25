package gohost

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ServeGRPC starts a gRPC endpoint for the given service.
func ServeGRPC(service GRPCService, grpcAddr string, opts []grpc.ServerOption) error {
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
	grpcServer := grpc.NewServer(opts...)
	service.RegisterServer(grpcServer)

	// start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC endpoint: %v", err)
	}

	return nil
}

// ServeGRPCWithTLS starts a gRPC endpoint for the given service with TLS enabled.
func ServeGRPCWithTLS(service GRPCService, grpcAddr string, opts []grpc.ServerOption, certFile string, keyFile string) error {
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
		return fmt.Errorf("failed to generate TLS credentials: %v", err)
	}

	// add TLS credentials to options
	opts = append(opts, grpc.Creds(creds))

	// start server
	return ServeGRPC(service, grpcAddr, opts)
}

// ServeHTTP starts an HTTP endpoint for a given service. This is a gateway pointing to a gRPC endpoint.
func ServeHTTP(service HTTPService, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption) error {
	// do not use TLS
	opts = append(opts, grpc.WithInsecure())

	// start server
	return serveHTTPInternal(service, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler http.Handler) error {
		return http.ListenAndServe(addr, handler)
	})
}

// ServeHTTPWithTLS starts an HTTP endpoint for a given service with TLS enabled. This is a gateway pointing to a gRPC endpoint.
func ServeHTTPWithTLS(service HTTPService, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption, certFile string, keyFile string, insecureSkipVerify bool) error {
	// validate parameters
	if certFile == "" {
		return errors.New("cert file cannot be empty")
	}
	if keyFile == "" {
		return errors.New("key file cannot be empty")
	}

	// create TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	})

	// add TLS credentials to options
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// start server
	return serveHTTPInternal(service, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler http.Handler) error {
		return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	})
}

// serveHTTPInternal is an internal method for serving up an HTTP endpoint.
func serveHTTPInternal(service HTTPService, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption, listenAndServe func(addr string, handler http.Handler) error) error {
	// validate parameters
	if service == nil {
		return errors.New("service cannot be nil")
	}
	if httpAddr == "" {
		return errors.New("http address cannot be empty")
	}
	if grpcAddr == "" {
		return errors.New("grpc address cannot be empty")
	}

	// create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// register server
	mux := runtime.NewServeMux()
	err := service.RegisterHandler(ctx, mux, grpcAddr, opts)
	if err != nil {
		return fmt.Errorf("failed to register HTTP endpoint: %v", err)
	}

	// enable CORS if requested
	var handler http.Handler = mux
	if enableCORS {
		handler = cors.AllowAll().Handler(mux)
	}

	// start server
	err = listenAndServe(httpAddr, handler)
	if err != nil {
		return fmt.Errorf("failed to serve HTTP endpoint: %v", err)
	}

	return nil
}
