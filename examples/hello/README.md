# Example Service: hello

A simple example service that demonstrates how to use a Hoster instance to host a service with gRPC and HTTP endpoints.

## Prerequisites
- Install [gRPC](https://grpc.io/docs/quickstart/go.html)
    - Make sure protoc is in GOPATH/bin
    - Make sure google/protobuf is also in GOPATH/bin
- Install [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

## Run the server
- Insecure
    - `go run cmd/server/main.go`
- With TLS
    - `go run cmd/server/main.go -cert-file ../../test/testdata/test.crt -key-file ../../test/testdata/test.key -insecure-skip-verify`

NOTE: insecure-skip-verify is only used for testing when the host name does not need to be verified and should not be used in production.

## Test the gRPC endpoint with the command-line client
- Insecure
    - `go run cmd/client/main.go -insecure -name eleniums`
- With TLS
    - `go run cmd/client/main.go -insecure-skip-verify -name eleniums`

## Test the HTTP endpoint with curl
- Insecure
    - `curl http://127.0.0.1:9090/v1/hello?name=eleniums`
- With TLS
    - `curl -k https://127.0.0.1:9090/v1/hello?name=eleniums`

## Test the debug endpoint
- Enable the debug endpoint when running the service:
    - `go run cmd/server/main.go -enable-debug`
- Navigate to:
    - http://127.0.0.1:6060/debug/pprof
    - http://127.0.0.1:6060/debug/vars

## Regenerate client/server from proto
- Use go:generate to build client/server/swagger:
    - `go generate`
