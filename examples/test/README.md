# Example Service: test

This example service is used for testing purposes.

## Prerequisites
- Install [gRPC](https://grpc.io/docs/quickstart/go.html)
    - Make sure protoc is in GOPATH/bin
    - Make sure google/protobuf is also in GOPATH/bin
- Install [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

## Regenerate client/server from proto
- Use go:generate to build client/server and swagger docs:
    - `go generate`
