package main

import (
	"crypto/tls"
	"flag"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/eleniums/gohost/examples/hello/proto"
)

func main() {
	// command-line flags
	grpcAddr := flag.String("grpc-addr", "127.0.0.1:50051", "host and port to host the gRPC endpoint")
	name := flag.String("name", "eleniums", "name to use in server greeting")
	insecure := flag.Bool("insecure", false, "true to use insecure connection and disable TLS")
	insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "true to skip verifying the certificate chain and host name")
	flag.Parse()

	// determine transport security to use
	var creds grpc.DialOption
	if *insecure {
		creds = grpc.WithInsecure()
	} else {
		creds = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: *insecureSkipVerify,
		}))
	}

	// dial the service
	conn, err := grpc.Dial(*grpcAddr, creds)
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
