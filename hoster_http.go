package gohost

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	h "net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	g "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// serveHTTP will start the HTTP endpoint.
func (h *Hoster) serveHTTP() error {
	// check if HTTP endpoint is enabled
	if h.HTTPAddr != "" {
		// ensure interface is implemented
		httpService, ok := h.Service.(gohttp.HTTPService)
		if !ok {
			return errors.New("service does not implement HTTP interface")
		}

		// configure dial options
		dialOpts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32), grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		}

		// start the HTTP endpoint
		if h.IsTLSEnabled() {
			go func() {
				gohttp.ServeHTTPWithTLS(httpService, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts, h.CertFile, h.KeyFile, h.InsecureSkipVerify)
			}()
		} else {
			go func() {
				gohttp.ServeHTTP(httpService, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts)
			}()
		}
	}

	return nil
}

// ServeHTTP starts an HTTP endpoint for a given service. This is a gateway pointing to a gRPC endpoint.
func ServeHTTP(service HTTPService, httpAddr string, grpcAddr string, enableCORS bool, opts []g.DialOption) error {
	// do not use TLS
	opts = append(opts, g.WithInsecure())

	// start server
	return serveHTTPInternal(service, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler h.Handler) error {
		return h.ListenAndServe(addr, handler)
	})
}

// ServeHTTPWithTLS starts an HTTP endpoint for a given service with TLS enabled. This is a gateway pointing to a gRPC endpoint.
func ServeHTTPWithTLS(service HTTPService, httpAddr string, grpcAddr string, enableCORS bool, opts []g.DialOption, certFile string, keyFile string, insecureSkipVerify bool) error {
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
	opts = append(opts, g.WithTransportCredentials(creds))

	// start server
	return serveHTTPInternal(service, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler h.Handler) error {
		return h.ListenAndServeTLS(addr, certFile, keyFile, handler)
	})
}

// serveHTTPInternal is an internal method for serving up an HTTP endpoint.
func serveHTTPInternal(service HTTPService, httpAddr string, grpcAddr string, enableCORS bool, opts []g.DialOption, listenAndServe func(addr string, handler h.Handler) error) error {
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
		return fmt.Errorf("failed to register HTTP handler: %v", err)
	}

	// enable CORS if requested
	var handler h.Handler = mux
	if enableCORS {
		handler = cors.AllowAll().Handler(mux)
	}

	// start server
	return listenAndServe(httpAddr, handler)
}
