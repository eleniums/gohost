package test

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/test/proto"
)

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

// RegisterServer registers the gRPC server to use with a service.
func (s *Service) RegisterServer(grpc *grpc.Server) {
	pb.RegisterTestServiceServer(grpc, s)
}

// RegisterHandler registers the HTTP handler to use with a service.
func (s *Service) RegisterHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return pb.RegisterTestServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
