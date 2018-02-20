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

// serveHTTP will start the HTTP endpoint.
func (h *Hoster) serveHTTP() error {
	// validate parameters
	if h.HTTPAddr == "" {
		return errors.New("http address cannot be empty")
	}

	// configure dial options
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32), grpc.MaxCallRecvMsgSize(math.MaxInt32)),
	}

	if h.isTLSEnabled() {
		// add TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: h.InsecureSkipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// add insecure option
		opts = append(opts, grpc.WithInsecure())
	}

	// create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// register servers
	mux := runtime.NewServeMux()
	for i := range h.httpHandlers {
		err := h.httpHandlers[i](ctx, mux, h.GRPCAddr, opts)
		if err != nil {
			return fmt.Errorf("failed to register HTTP handler: %v", err)
		}
	}

	// enable CORS if requested
	var handler http.Handler = mux
	if h.EnableCORS {
		handler = cors.AllowAll().Handler(mux)
	}

	// start the HTTP endpoint
	if h.isTLSEnabled() {
		return http.ListenAndServeTLS(h.HTTPAddr, h.CertFile, h.KeyFile, handler)
	}

	return http.ListenAndServe(h.HTTPAddr, handler)
}
