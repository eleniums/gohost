package hello

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/hello/proto"
)

type Server struct {
}

// Hello will return a personalized greeting.
func (s *Server) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Greeting: fmt.Sprintf("Hello %v!", in.Name),
	}, nil
}

// RegisterServer registers this server to be a gRPC endpoint.
func (s *Server) RegisterServer(grpc *grpc.Server) {
	pb.RegisterHelloServiceServer(grpc, s)
}

// RegisterHandler registers this server to be an HTTP endpoint.
func (s *Server) RegisterHandler(ctx context.Context, handler http.Handler, endpoint string, opts []grpc.DialOption) error {
	return pb.RegisterHelloServiceHandlerFromEndpoint(ctx, handler, endpoint, opts)
}
