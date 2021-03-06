package main

import (
	"flag"
	"log"

	"github.com/eleniums/gohost"
	"github.com/eleniums/gohost/examples/hello"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/hello/proto"
)

func main() {
	// command-line flags
	grpcAddr := flag.String("grpc-addr", "127.0.0.1:50051", "host and port to host the gRPC endpoint")
	httpAddr := flag.String("http-addr", "127.0.0.1:9090", "host and port to host the HTTP endpoint")
	debugAddr := flag.String("debug-addr", "127.0.0.1:6060", "host and port to host the debug endpoint (/debug/pprof and /debug/vars)")
	enableDebug := flag.Bool("enable-debug", false, "true to enable the debug endpoint (/debug/pprof and /debug/vars)")
	certFile := flag.String("cert-file", "", "cert file for enabling a TLS connection")
	keyFile := flag.String("key-file", "", "key file for enabling a TLS connection")
	insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "true to skip verifying the certificate chain and host name")
	maxSendMsgSize := flag.Int("max-send-msg-size", gohost.DefaultMaxSendMsgSize, "max message size the service is allowed to send")
	maxRecvMsgSize := flag.Int("max-recv-msg-size", gohost.DefaultMaxRecvMsgSize, "max message size the service is allowed to receive")
	flag.Parse()

	// create the service
	service := hello.NewService()

	// create the hoster
	hoster := gohost.NewHoster()
	hoster.GRPCAddr = *grpcAddr
	hoster.HTTPAddr = *httpAddr
	hoster.DebugAddr = *debugAddr
	hoster.EnableDebug = *enableDebug
	hoster.CertFile = *certFile
	hoster.KeyFile = *keyFile
	hoster.InsecureSkipVerify = *insecureSkipVerify
	hoster.MaxSendMsgSize = *maxSendMsgSize
	hoster.MaxRecvMsgSize = *maxRecvMsgSize

	hoster.RegisterGRPCServer(func(s *grpc.Server) {
		pb.RegisterHelloServiceServer(s, service)
	})
	log.Printf("Registered gRPC endpoint at: %v", *grpcAddr)

	hoster.RegisterHTTPGateway(pb.RegisterHelloServiceHandlerFromEndpoint)
	log.Printf("Registered HTTP endpoint at: %v", *httpAddr)

	// start the server
	err := hoster.ListenAndServe()
	if err != nil {
		log.Fatalf("Unable to start the server: %v", err)
	}
}
