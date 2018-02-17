package gohost

import (
	"context"
	"errors"
	"math"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	// register debug http handlers
	_ "expvar"
	_ "net/http/pprof"
)

const (
	// DefaultMaxSendMsgSize is the default max send message size, per gRPC
	DefaultMaxSendMsgSize = 1024 * 1024 * 4

	// DefaultMaxRecvMsgSize is the default max receive message size, per gRPC
	DefaultMaxRecvMsgSize = 1024 * 1024 * 4
)

// Hoster is used to serve gRPC and HTTP endpoints.
type Hoster struct {
	// GRPCAddr is the endpoint (host and port) on which to host the gRPC services.
	GRPCAddr string

	// HTTPAddr is the endpoint (host and port) on which to host the HTTP services. May be left blank if not using HTTP.
	HTTPAddr string

	// DebugAddr is the endpoint (host and port) on which to host the debug endpoint (/debug/pprof and /debug/vars). May be left blank to disable debug endpoint.
	DebugAddr string

	// CertFile is the certificate file for use with TLS. May be left blank if using insecure mode.
	CertFile string

	// KeyFile is the private key file for use with TLS. May be left blank if using insecure mode.
	KeyFile string

	// InsecureSkipVerify will cause verification of the host name during a TLS handshake to be skipped if set to true.
	InsecureSkipVerify bool

	// EnableCORS will enable all cross-origin resource sharing if set to true.
	EnableCORS bool

	// MaxSendMsgSize will change the size of the message that can be sent from the service.
	MaxSendMsgSize int

	// MaxRecvMsgSize will change the size of the message that can be received by the service.
	MaxRecvMsgSize int

	// UnaryInterceptors is an array of unary interceptors to be used by the service. They will be executed in order, from first to last.
	UnaryInterceptors []grpc.UnaryServerInterceptor

	// StreamInterceptors is an array of stream interceptors to be used by the service. They will be executed in order, from first to last.
	StreamInterceptors []grpc.StreamServerInterceptor

	grpcEndpoints []func(s *grpc.Server)
	httpEndpoints []func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

// NewHoster creates a new hoster instance with defaults set. This is the minimum required to host a server.
func NewHoster() *Hoster {
	return &Hoster{
		MaxSendMsgSize: DefaultMaxSendMsgSize,
		MaxRecvMsgSize: DefaultMaxRecvMsgSize,
	}
}

func (h *Hoster) AddGRPCEndpoint(endpoint ...func(s *grpc.Server)) {
	h.grpcEndpoints = append(h.grpcEndpoints, endpoint...)
}

func (h *Hoster) AddHTTPGateway(gateway ...func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)) {
	h.httpEndpoints = append(h.httpEndpoints, gateway...)
}

// ListenAndServe creates and starts the server.
func (h *Hoster) ListenAndServe() error {
	// validate parameters
	if h.Service == nil {
		return errors.New("gRPC service implementation must be provided")
	}
	if h.GRPCAddr == "" {
		return errors.New("gRPC address must be provided")
	}

	// serve debug endpoint
	h.serveDebug()

	// serve HTTP endpoint
	err := h.serveHTTP()
	if err != nil {
		return err
	}

	// serve gRPC endpoint
	return h.serveGRPC()
}

// IsTLSEnabled will return true if TLS properties are set and ready to use.
func (h *Hoster) IsTLSEnabled() bool {
	return h.CertFile != "" && h.KeyFile != ""
}

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

// serveDebug will start the debug endpoint.
func (h *Hoster) serveDebug() {
	// check if debug endpoint is enabled
	if h.DebugAddr != "" {
		go func() {
			http.ListenAndServe(h.DebugAddr, nil)
		}()
	}
}
