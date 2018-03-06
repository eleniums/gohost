package test

import (
	"io"

	"golang.org/x/net/context"

	pb "github.com/eleniums/gohost/examples/test/proto"
)

//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. proto/test.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. proto/test.proto
//go:generate protoc -I. -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --proto_path=./proto --swagger_out=logtostderr=true:. proto/test.proto

// Service contains the implementation for the gRPC service. It implements the interfaces for both gRPC and HTTP endpoints.
type Service struct{}

// NewService creates a new instance of Service.
func NewService() *Service {
	return &Service{}
}

// Send the value in the request.
func (s *Service) Send(ctx context.Context, in *pb.SendRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{
		Success: true,
	}, nil
}

// Echo the value in the request back in the response.
func (s *Service) Echo(ctx context.Context, in *pb.SendRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Echo: in.Value,
	}, nil
}

// Large will send a large response message.
func (s *Service) Large(ctx context.Context, in *pb.LargeRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Echo: string(make([]byte, in.Length)),
	}, nil
}

// Stream a bunch of requests.
func (s *Service) Stream(stream pb.TestService_StreamServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.TestResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}
	}
}
