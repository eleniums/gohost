# Example Service: hello

A simple example service that demonstrates how to use a Hoster instance to host a service with gRPC and HTTP endpoints.

## Prerequisites
- Install [gRPC](https://grpc.io/docs/quickstart/go.html)
    - Make sure protoc is in GOPATH/bin
    - Make sure google/protobuf is also in GOPATH/bin
- Install [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
    - Install from vendor directory to avoid issues
        - https://github.com/grpc-ecosystem/grpc-gateway/issues/384#issuecomment-300863457

## Run the server
- Insecure
    - `go run cmd/server/main.go`
- With TLS
    - `go run cmd/server/main.go -cert-file ../../testdata/test.crt -key-file ../../testdata/test.key -insecure-skip-verify`

NOTE: insecure-skip-verify is only used for testing when the host name does not need to be verified and should not be used in production.

## Test the gRPC endpoint with the command-line client
- Insecure
    - `go run cmd/cli/main.go -insecure -name eleniums`
- With TLS
    - `go run cmd/cli/main.go -insecure-skip-verify -name eleniums`

## Test the HTTP endpoint with curl
- Insecure
    - `curl http://127.0.0.1:9090/v1/hello?name=eleniums`
- With TLS
    - `curl -k https://127.0.0.1:9090/v1/hello?name=eleniums`

## Test the debug endpoint
- Navigate to:
    - http://127.0.0.1:6060/debug/pprof
    - http://127.0.0.1:6060/debug/vars

## Generate Stubs

The stubs for this example have already been generated and checked in, so these commands are only provided as a reference.

- Generate gRPC client/server stubs:
    - `protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./ proto/hello.proto`
- Generate HTTP gateway:
    - `protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. proto/hello.proto`
- Generate Swagger definitions:
    - `protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --proto_path=./proto --swagger_out=logtostderr=true:. proto/hello.proto`
