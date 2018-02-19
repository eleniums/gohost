package gohost

import (
	"google.golang.org/grpc"
)

const (
	// DefaultDebugAddr is the default address for the debug endpoint (/debug/pprof and /debug/vars).
	DefaultDebugAddr = "127.0.0.1:6060"

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

	// DebugAddr is the endpoint (host and port) on which to host the debug endpoint (/debug/pprof and /debug/vars). Default is 127.0.0.1:6060.
	DebugAddr string

	// CertFile is the certificate file for use with TLS. May be left blank if using insecure mode.
	CertFile string

	// KeyFile is the private key file for use with TLS. May be left blank if using insecure mode.
	KeyFile string

	// InsecureSkipVerify will cause verification of the host name during a TLS handshake to be skipped if set to true.
	InsecureSkipVerify bool

	// EnableCORS will enable all cross-origin resource sharing if set to true.
	EnableCORS bool

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

	grpcEndpoints []grpcServer
	httpEndpoints []httpGateway
}

// NewHoster creates a new hoster instance with defaults set. This is the minimum required to host a server.
func NewHoster() *Hoster {
	return &Hoster{
		DebugAddr:      DefaultDebugAddr,
		MaxSendMsgSize: DefaultMaxSendMsgSize,
		MaxRecvMsgSize: DefaultMaxRecvMsgSize,
	}
}

func (h *Hoster) AddGRPCEndpoint(endpoint ...grpcServer) {
	h.grpcEndpoints = append(h.grpcEndpoints, endpoint...)
}

func (h *Hoster) AddHTTPGateway(gateway ...httpGateway) {
	h.httpEndpoints = append(h.httpEndpoints, gateway...)
}

// ListenAndServe creates and starts the server.
func (h *Hoster) ListenAndServe() error {
	// serve debug endpoint
	if h.EnableDebug {
		go func() {
			h.serveDebug()
		}()
	}

	// serve HTTP endpoint
	if len(h.httpEndpoints) > 0 {
		go func() {
			h.serveHTTP()
		}()
		// if err != nil {
		// 	return err
		// }
	}

	// serve gRPC endpoint
	return h.serveGRPC()
}

// IsTLSEnabled will return true if TLS properties are set and ready to use.
func (h *Hoster) IsTLSEnabled() bool {
	return h.CertFile != "" && h.KeyFile != ""
}
