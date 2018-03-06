package gohost

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// DefaultGRPCAddr is the default address for the gRPC endpoint.
	DefaultGRPCAddr = "127.0.0.1:50051"

	// DefaultHTTPAddr is the default address for the HTTP endpoint.
	DefaultHTTPAddr = "127.0.0.1:9090"

	// DefaultDebugAddr is the default address for the debug endpoint (/debug/pprof and /debug/vars).
	DefaultDebugAddr = "127.0.0.1:6060"

	// DefaultMaxSendMsgSize is the default max send message size, per gRPC
	DefaultMaxSendMsgSize = 1024 * 1024 * 4

	// DefaultMaxRecvMsgSize is the default max receive message size, per gRPC
	DefaultMaxRecvMsgSize = 1024 * 1024 * 4
)

// GRPCEndpoint is used to register a gRPC endpoint.
type GRPCEndpoint func(s *grpc.Server)

// HTTPEndpoint is used to register a HTTP endpoint for forwarding requests to a gRPC endpoint.
type HTTPEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// Hoster is used to serve gRPC and HTTP endpoints.
type Hoster struct {
	// GRPCAddr is the endpoint (host and port) on which to host the gRPC services. Default is 127.0.0.1:50051. May be left blank if no gRPC endpoints have been registered.
	GRPCAddr string

	// HTTPAddr is the endpoint (host and port) on which to host the HTTP services. Default is 127.0.0.1:9090. May be left blank if no HTTP endpoints have been registered.
	HTTPAddr string

	// DebugAddr is the endpoint (host and port) on which to host the debug endpoint (/debug/pprof and /debug/vars). Default is 127.0.0.1:6060. May be left blank if EnableDebug is false.
	DebugAddr string

	// CertFile is the certificate file for use with TLS. May be left blank if using insecure mode.
	CertFile string

	// KeyFile is the private key file for use with TLS. May be left blank if using insecure mode.
	KeyFile string

	// InsecureSkipVerify will cause verification of the host name during a TLS handshake to be skipped if set to true.
	InsecureSkipVerify bool

	// HTTPHandler is used to register a handler that can optionally be added to the HTTP endpoint. Leave blank to use default mux.
	HTTPHandler func(mux *runtime.ServeMux) http.Handler

	// EnableDebug will enable the debug endpoint (/debug/pprof and /debug/vars). The debug endpoint address is defined by DebugAddr.
	EnableDebug bool

	// MaxSendMsgSize will change the size of the message that can be sent from the service.
	MaxSendMsgSize int

	// MaxRecvMsgSize will change the size of the message that can be received by the service.
	MaxRecvMsgSize int

	// UnaryInterceptors is an array of unary interceptors to be used by the service. They will be executed in order, from first to last.
	UnaryInterceptors []grpc.UnaryServerInterceptor

	// StreamInterceptors is an array of stream interceptors to be used by the service. They will be executed in order, from first to last.
	StreamInterceptors []grpc.StreamServerInterceptor

	// grpcEndpoints is an array of gRPC endpoints to be hosted.
	grpcEndpoints []GRPCEndpoint

	// httpEndpoints is an array of HTTP endpoints to be hosted.
	httpEndpoints []HTTPEndpoint
}

// NewHoster creates a new hoster instance with defaults set.
func NewHoster() *Hoster {
	return &Hoster{
		GRPCAddr:       DefaultGRPCAddr,
		HTTPAddr:       DefaultHTTPAddr,
		DebugAddr:      DefaultDebugAddr,
		MaxSendMsgSize: DefaultMaxSendMsgSize,
		MaxRecvMsgSize: DefaultMaxRecvMsgSize,
	}
}

// RegisterGRPCEndpoint will add a function for registering a gRPC endpoint. The function is invoked when ListenAndServe is called.
func (h *Hoster) RegisterGRPCEndpoint(endpoints ...GRPCEndpoint) {
	h.grpcEndpoints = append(h.grpcEndpoints, endpoints...)
}

// RegisterHTTPEndpoint will add a function for registering a HTTP endpoint. The function is invoked when ListenAndServe is called.
func (h *Hoster) RegisterHTTPEndpoint(handlers ...HTTPEndpoint) {
	h.httpEndpoints = append(h.httpEndpoints, handlers...)
}

// ListenAndServe creates and starts the server.
func (h *Hoster) ListenAndServe() error {
	errc := make(chan error)

	// serve debug endpoint
	if h.EnableDebug {
		go func() {
			errc <- h.serveDebug()
		}()
	}

	// serve HTTP endpoint
	if len(h.httpEndpoints) > 0 {
		go func() {
			errc <- h.serveHTTP()
		}()
	}

	// serve gRPC endpoint
	if len(h.grpcEndpoints) > 0 {
		go func() {
			errc <- h.serveGRPC()
		}()
	}

	return <-errc
}
