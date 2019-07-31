# gohost

[![Build Status](https://travis-ci.org/eleniums/gohost.svg?branch=master)](https://travis-ci.org/eleniums/gohost) [![Go Report Card](https://goreportcard.com/badge/github.com/eleniums/gohost)](https://goreportcard.com/report/github.com/eleniums/gohost) [![GoDoc](https://godoc.org/github.com/eleniums/gohost?status.svg)](https://godoc.org/github.com/eleniums/gohost) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/eleniums/gohost/blob/master/LICENSE)

A tool for hosting a Go service with gRPC and HTTP endpoints.

**It is generally better to just use the standard libraries directly. Less bloat and more control over configuration. See:**
- net/http: https://golang.org/pkg/net/http
- gRPC: https://github.com/grpc/grpc-go

## Installation

```
go get -u github.com/eleniums/gohost
```

## Prerequisites

- Requires Go 1.9 or later
- Uses [dep](https://github.com/golang/dep) for dependencies
- Uses [grpc-go](https://github.com/grpc/grpc-go) for gRPC endpoints
- Uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) for HTTP endpoints
- See the full list of imported packages [here](https://godoc.org/github.com/eleniums/gohost?imports)

## Example

Sample service implementation:
```go
// Service contains the implementation for the gRPC service.
type Service struct{}

// NewService creates a new instance of Service.
func NewService() *Service {
	return &Service{}
}

// Hello will return a personalized greeting.
func (s *Service) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// create greeting
	greeting := "Hello!"
	if in.Name != "" {
		greeting = fmt.Sprintf("Hello %v!", in.Name)
	}

	log.Printf("Received request from: %v", in.Name)

	// return response
	return &pb.HelloResponse{
		Greeting: greeting,
	}, nil
}
```

Use the Hoster struct to serve up gRPC and HTTP endpoints:
```go
// create the service
service := hello.NewService()

// create the hoster
hoster := gohost.NewHoster()
hoster.GRPCAddr = *grpcAddr
hoster.HTTPAddr = *httpAddr
hoster.DebugAddr = *debugAddr
hoster.EnableDebug = *enableDebug
hoster.CertFile = *certFile
hoster.KeyFile = *keyFile
hoster.InsecureSkipVerify = *insecureSkipVerify
hoster.MaxSendMsgSize = *maxSendMsgSize
hoster.MaxRecvMsgSize = *maxRecvMsgSize

hoster.RegisterGRPCServer(func(s *grpc.Server) {
	pb.RegisterHelloServiceServer(s, service)
})

hoster.RegisterHTTPGateway(pb.RegisterHelloServiceHandlerFromEndpoint)

// start the server
err := hoster.ListenAndServe()
if err != nil {
	log.Fatalf("Unable to start the server: %v", err)
}
```

See the full example [here](https://github.com/eleniums/gohost/tree/master/examples/hello).
