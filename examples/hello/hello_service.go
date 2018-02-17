package hello

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

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
