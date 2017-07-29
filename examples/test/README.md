# Example Service: test

# Prerequisites
- Install gRPC
    - https://grpc.io/docs/quickstart/go.html
    - Make sure protoc is in GOPATH/bin
    - Make sure google/protobuf is also in GOPATH/bin
- Install grpc-gateway
    - https://github.com/grpc-ecosystem/grpc-gateway
    - Install from vendor directory to avoid issues
        - https://github.com/grpc-ecosystem/grpc-gateway/issues/384#issuecomment-300863457

# Generate Endpoints
- Generate gRPC client/server
    - protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./ proto/test.proto
- Generate HTTP gateway
    - protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. proto/test.proto
- Generate Swagger output
    - protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --proto_path=./proto --swagger_out=logtostderr=true:. proto/test.proto