package gohost

import (
	"errors"

	"google.golang.org/grpc"
)

const (
	// DefaultMaxSendMsgSize is the default max send message size, per gRPC
	DefaultMaxSendMsgSize = 1024 * 1024 * 4

	// DefaultMaxRecvMsgSize is the default max receive message size, per gRPC
	DefaultMaxRecvMsgSize = 1024 * 1024 * 4
)

// Hoster is used to serve gRPC and HTTP endpoints.
type Hoster struct {
	Server             Server
	GRPCAddr           string
	HTTPAddr           string
	CertFile           string
	KeyFile            string
	InsecureSkipVerify bool
	EnableCORS         bool
	MaxSendMsgSize     int
	MaxRecvMsgSize     int
}

// NewHoster creates a new hoster instance with defaults set. This is the minimum required to host a server.
func NewHoster(server Server, grpcAddr string) *Hoster {
	return &Hoster{
		Server:         server,
		GRPCAddr:       grpcAddr,
		MaxSendMsgSize: DefaultMaxSendMsgSize,
		MaxRecvMsgSize: DefaultMaxRecvMsgSize,
	}
}

// ListenAndServe creates and starts the server.
func (h *Hoster) ListenAndServe() error {
	// validate parameters
	if h.GRPCAddr == "" {
		return errors.New("gRPC address must be provided")
	}

	// check if HTTP endpoint is enabled
	if h.HTTPAddr != "" {
		// configure dial options
		dialOpts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(h.MaxSendMsgSize), grpc.MaxCallRecvMsgSize(h.MaxRecvMsgSize)),
		}

		// start the HTTP endpoint
		if h.IsTLSEnabled() {
			go ServeHTTPWithTLS(h.Server, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts, h.CertFile, h.KeyFile, h.InsecureSkipVerify)
		} else {
			go ServeHTTP(h.Server, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts)
		}
	}

	// configure server options
	serverOpts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(h.MaxSendMsgSize),
		grpc.MaxRecvMsgSize(h.MaxRecvMsgSize),
	}

	// start the gRPC endpoint
	if h.IsTLSEnabled() {
		return ServeGRPCWithTLS(h.Server, h.GRPCAddr, serverOpts, h.CertFile, h.KeyFile)
	}

	return ServeGRPC(h.Server, h.GRPCAddr, serverOpts)
}

// IsTLSEnabled will return true if TLS properties are set and ready to use.
func (h *Hoster) IsTLSEnabled() bool {
	return h.CertFile != "" && h.KeyFile != ""
}
