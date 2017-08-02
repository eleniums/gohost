package hello

import (
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/hello/proto"
)

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
