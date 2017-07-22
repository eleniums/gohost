package main

import (
	"flag"
	"log"

	"github.com/eleniums/gohost"
	"github.com/eleniums/gohost/examples/hello"
)

func main() {
	// command-line flags
	grpcAddr := flag.String("grpc-addr", "127.0.0.1:50051", "host and port to host the gRPC endpoint")
	httpAddr := flag.String("http-addr", "127.0.0.1:9090", "host and port to host the HTTP endpoint")
	enableCors := flag.Bool("enable-cors", false, "true to enable cross-origin resource sharing (CORS)")
	certFile := flag.String("cert-file", "", "cert file for enabling a TLS connection")
	keyFile := flag.String("key-file", "", "key file for enabling a TLS connection")
	insecureSkipVerify := flag.Bool("insecure-skip-verify", false, "true to skip verifying the certificate chain and host name")
	maxSendMsgSize := flag.Int("max-send-msg-size", gohost.DefaultMaxSendMsgSize, "max message size the service is allowed to send")
	maxRecvMsgSize := flag.Int("max-recv-msg-size", gohost.DefaultMaxRecvMsgSize, "max message size the service is allowed to receive")
	flag.Parse()

	// create the service
	service := hello.NewService()

	// create the hoster
	hoster := gohost.NewHoster(service, *grpcAddr)
	hoster.HTTPAddr = *httpAddr
	hoster.EnableCORS = *enableCors
	hoster.CertFile = *certFile
	hoster.KeyFile = *keyFile
	hoster.InsecureSkipVerify = *insecureSkipVerify
	hoster.MaxSendMsgSize = *maxSendMsgSize
	hoster.MaxRecvMsgSize = *maxRecvMsgSize
	hoster.Logger = log.Printf

	// start the server
	err := hoster.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start the server: %v", err)
	}
}
