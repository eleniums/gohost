package grpc

// GRPCEndpoint contains the properties necessary to host a gRPC endpoint.
type GRPCEndpoint struct {
}

// Serve hosts the gRPC endpoint.
func (e *GRPCEndpoint) Serve() error {
	return nil
}
