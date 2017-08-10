package test

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/test/proto"
)

// GRPCService contains the implementation for the gRPC service. It does not implement the interface for the HTTP endpoint.
type GRPCService struct{}

// NewGRPCService creates a new instance of GRPCService.
func NewGRPCService() *GRPCService {
	return &GRPCService{}
}

// Send the value in the request.
func (s *GRPCService) Send(ctx context.Context, in *pb.SendRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{
		Success: true,
	}, nil
}

// Echo the value in the request back in the response.
func (s *GRPCService) Echo(ctx context.Context, in *pb.SendRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Echo: in.Value,
	}, nil
}

// Large will send a large response message.
func (s *GRPCService) Large(ctx context.Context, in *pb.LargeRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{
		Echo: string(make([]byte, in.Length)),
	}, nil
}

// RegisterServer registers the gRPC server to use with a service.
func (s *GRPCService) RegisterServer(grpc *grpc.Server) {
	pb.RegisterTestServiceServer(grpc, s)
}
