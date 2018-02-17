package gohost

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type httpGateway func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// serveHTTP will start the HTTP endpoint.
func (h *Hoster) serveHTTP() error {
	// validate parameters
	if len(h.httpEndpoints) == 0 {
		return errors.New("no http gateways added")
	}
	if h.HTTPAddr == "" {
		return errors.New("http address cannot be empty")
	}

	// check if HTTP endpoint is enabled
	if h.HTTPAddr != "" {
		// configure dial options
		dialOpts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32), grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		}

		// start the HTTP endpoint
		if h.IsTLSEnabled() {
			go func() {
				ServeHTTPWithTLS(h.httpEndpoints, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts, h.CertFile, h.KeyFile, h.InsecureSkipVerify)
			}()
		} else {
			go func() {
				ServeHTTP(h.httpEndpoints, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts)
			}()
		}
	}

	return nil
}

// ServeHTTP starts an HTTP endpoint for a given service. This is a gateway pointing to a gRPC endpoint.
func ServeHTTP(gateways []httpGateway, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption) error {
	// do not use TLS
	opts = append(opts, grpc.WithInsecure())

	// start server
	return serveHTTPInternal(gateways, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler http.Handler) error {
		return http.ListenAndServe(addr, handler)
	})
}

// ServeHTTPWithTLS starts an HTTP endpoint for a given service with TLS enabled. This is a gateway pointing to a gRPC endpoint.
func ServeHTTPWithTLS(gateways []httpGateway, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption, certFile string, keyFile string, insecureSkipVerify bool) error {
	// create TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	})

	// add TLS credentials to options
	opts = append(opts, grpc.WithTransportCredentials(creds))

	// start server
	return serveHTTPInternal(gateways, httpAddr, grpcAddr, enableCORS, opts, func(addr string, handler http.Handler) error {
		return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	})
}

// serveHTTPInternal is an internal method for serving up an HTTP endpoint.
func serveHTTPInternal(gateways []httpGateway, httpAddr string, grpcAddr string, enableCORS bool, opts []grpc.DialOption, listenAndServe func(addr string, handler http.Handler) error) error {
	// create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// register server
	mux := runtime.NewServeMux()
	for i := range gateways {
		err := gateways[i](ctx, mux, grpcAddr, opts)
		if err != nil {
			return fmt.Errorf("failed to register HTTP handler: %v", err)
		}
	}

	// enable CORS if requested
	var handler http.Handler = mux
	if enableCORS {
		handler = cors.AllowAll().Handler(mux)
	}

	// start server
	return listenAndServe(httpAddr, handler)
}
