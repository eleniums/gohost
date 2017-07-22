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
	Service            GRPCService
	GRPCAddr           string
	HTTPAddr           string
	CertFile           string
	KeyFile            string
	InsecureSkipVerify bool
	EnableCORS         bool
	MaxSendMsgSize     int
	MaxRecvMsgSize     int
	Log                func(format string, v ...interface{})
}

// NewHoster creates a new hoster instance with defaults set. This is the minimum required to host a server.
func NewHoster(service GRPCService, grpcAddr string) *Hoster {
	return &Hoster{
		Service:        service,
		GRPCAddr:       grpcAddr,
		MaxSendMsgSize: DefaultMaxSendMsgSize,
		MaxRecvMsgSize: DefaultMaxRecvMsgSize,
	}
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

	// check if HTTP endpoint is enabled
	if h.HTTPAddr != "" {
		// ensure interface is implemented
		httpService, ok := h.Service.(HTTPService)
		if !ok {
			return errors.New("service does not implement HTTP interface")
		}

		// configure dial options
		dialOpts := []grpc.DialOption{
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(h.MaxSendMsgSize), grpc.MaxCallRecvMsgSize(h.MaxRecvMsgSize)),
		}

		// start the HTTP endpoint
		if h.IsTLSEnabled() {
			h.Log("Starting HTTP endpoint with TLS enabled: %v", h.HTTPAddr)
			go ServeHTTPWithTLS(httpService, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts, h.CertFile, h.KeyFile, h.InsecureSkipVerify)
		} else {
			h.Log("Starting insecure HTTP endpoint: %v", h.HTTPAddr)
			go ServeHTTP(httpService, h.HTTPAddr, h.GRPCAddr, h.EnableCORS, dialOpts)
		}
	}

	// configure server options
	serverOpts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(h.MaxSendMsgSize),
		grpc.MaxRecvMsgSize(h.MaxRecvMsgSize),
	}

	// start the gRPC endpoint
	if h.IsTLSEnabled() {
		h.Log("Starting gRPC endpoint with TLS enabled: %v", h.GRPCAddr)
		return ServeGRPCWithTLS(h.Service, h.GRPCAddr, serverOpts, h.CertFile, h.KeyFile)
	}

	h.Log("Starting insecure gRPC endpoint: %v", h.GRPCAddr)
	return ServeGRPC(h.Service, h.GRPCAddr, serverOpts)
}

// IsTLSEnabled will return true if TLS properties are set and ready to use.
func (h *Hoster) IsTLSEnabled() bool {
	return h.CertFile != "" && h.KeyFile != ""
}
