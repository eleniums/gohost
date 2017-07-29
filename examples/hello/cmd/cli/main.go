package main

import (
	"flag"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/hello/proto"
)

func main() {
	// command-line flags
	grpcAddr := flag.String("grpc-addr", "127.0.0.1:50051", "host and port to host the gRPC endpoint")
	name := flag.String("name", "eleniums", "name to use in server greeting")
	flag.Parse()

	// dial the service
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dail service: %v", err)
	}

	// create client for service
	client := pb.NewHelloServiceClient(conn)

	// create a request
	request := pb.HelloRequest{
		Name: *name,
	}

	start := time.Now()

	// call the hello function on the server
	response, err := client.Hello(context.Background(), &request)
	if err != nil {
		log.Fatalf("failed to say hello: %v", err)
	}

	// display server response
	log.Printf("Server response in %v: %v", time.Since(start), response.Greeting)
}
