# gohost

[![Go Report Card](https://goreportcard.com/badge/github.com/eleniums/gohost)](https://goreportcard.com/report/github.com/eleniums/gohost) [![Coverage](http://gocover.io/_badge/github.com/eleniums/gohost)](http://gocover.io/github.com/eleniums/gohost) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/eleniums/gohost/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/eleniums/gohost?status.svg)](https://godoc.org/github.com/eleniums/gohost)

A tool for hosting a Go service with gRPC and HTTP endpoints.

## Installation

```
go get github.com/eleniums/gohost
```

## Prerequisites

- Requires Go 1.6 or later
- Uses [dep](https://github.com/golang/dep) for dependencies
- Uses [grpc-go](https://github.com/grpc/grpc-go) for gRPC endpoints
- Uses [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) for HTTP endpoints
- See the full list of imported packages [here](https://godoc.org/github.com/eleniums/gohost?imports)

## Example

Service implementing the [grpc_service](https://github.com/eleniums/gohost/blob/master/grpc_service.go) and [http_service](https://github.com/eleniums/gohost/blob/master/http_service.go) interfaces:
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

// RegisterServer registers the gRPC server to use with a service.
func (s *Service) RegisterServer(grpc *grpc.Server) {
	pb.RegisterHelloServiceServer(grpc, s)
}

// RegisterHandler registers the HTTP handler to use with a service.
func (s *Service) RegisterHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
```

Use the Hoster struct to serve up gRPC and HTTP endpoints:
```go
// create the service
service := hello.NewService()

// create the hoster
hoster := gohost.NewHoster(service, *grpcAddr)
hoster.HTTPAddr = *httpAddr
hoster.EnableCORS = *enableCors
hoster.CertFile = *certFile
hoster.KeyFile = *keyFile
hoster.InsecureSkipVerify = *insecureSkipVerify
hoster.MaxSendMsgSize = *maxSendMsgSize
hoster.MaxRecvMsgSize = *maxRecvMsgSize
hoster.Logger = log.Printf

// start the server
err := hoster.ListenAndServe()
if err != nil {
	log.Fatalf("unable to start the server: %v", err)
}
```

See the full example [here](https://github.com/eleniums/gohost/tree/master/examples/hello).

## Issues
- There is an issue with the latest version of grpc-gateway
    - https://github.com/grpc-ecosystem/grpc-gateway/issues/384#issuecomment-300863457
    - Install version 1.2.2 from the vendor directory
        - `go install ./github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
        - `go install ./github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger`
