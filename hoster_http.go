package gohost

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(h.MaxSendMsgSize), grpc.MaxCallRecvMsgSize(h.MaxRecvMsgSize)),
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

	// register gateways
	mux := runtime.NewServeMux()
	for i := range h.httpGateways {
		err := h.httpGateways[i](ctx, mux, h.GRPCAddr, opts)
		if err != nil {
			return fmt.Errorf("failed to register HTTP gateway: %v", err)
		}
	}

	// register optional handler
	var handler http.Handler = mux
	if h.HTTPHandler != nil {
		handler = h.HTTPHandler(mux)
	}

	// start the HTTP endpoint
	if h.isTLSEnabled() {
		return http.ListenAndServeTLS(h.HTTPAddr, h.CertFile, h.KeyFile, handler)
	}

	return http.ListenAndServe(h.HTTPAddr, handler)
}
