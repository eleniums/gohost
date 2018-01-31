package gohost

// HosterEndpoint is an interface for an endpoint that can be hosted.
type HosterEndpoint interface {
	// Serve hosts the endpoint.
	Serve() error
}
