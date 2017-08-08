# Example Service: test

This example service is used for testing purposes.

## Prerequisites
- Install [gRPC](https://grpc.io/docs/quickstart/go.html)
    - Make sure protoc is in GOPATH/bin
    - Make sure google/protobuf is also in GOPATH/bin
- Install [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
    - Install from vendor directory to avoid issues
        - https://github.com/grpc-ecosystem/grpc-gateway/issues/384#issuecomment-300863457

## Generate Stubs

The stubs for this example have already been generated and checked in, so these commands are only provided as a reference.

- Generate gRPC client/server stubs:
    - `protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./ proto/test.proto`
- Generate HTTP gateway:
    - `protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. proto/test.proto`
- Generate Swagger definitions:
    - `protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --proto_path=./proto --swagger_out=logtostderr=true:. proto/test.proto`
