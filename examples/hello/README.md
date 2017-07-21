# Example Service: hello

# Generate endpoints
- Install gRPC
    - https://grpc.io/docs/quickstart/go.html
    - Make sure protoc is in GOPATH/bin
    - Make sure google/protobuf is also in GOPATH/bin
- Generate gRPC client/server
    - protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./ proto/hello.proto
- Generate HTTP gateway
    - protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. proto/hello.proto
- Generate Swagger output
    - protoc -I ./ -I ../../../../../github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --proto_path=./proto --swagger_out=logtostderr=true:. proto/hello.proto

# TODO
- simplify/fix instructions for installing grpc
- need to change context to golang.org/x/net/context?