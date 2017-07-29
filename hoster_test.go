package gohost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/eleniums/gohost/examples/hello"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/hello/proto"
	assert "github.com/stretchr/testify/require"
)

func Test_Hoster_ListenAndServe_GRPCEndpoint(t *testing.T) {
	// arrange
	service := hello.NewService()
	grpcAddr := "127.0.0.1:50051"

	hoster := NewHoster(service, grpcAddr)

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewHelloServiceClient(conn)
	grpcReq := pb.HelloRequest{
		Name: "eleniums",
	}
	grpcResp, err := client.Hello(context.Background(), &grpcReq)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, "Hello eleniums!", grpcResp.Greeting)
}

func Test_Hoster_ListenAndServe_HTTPEndpoint(t *testing.T) {
	// arrange
	service := hello.NewService()
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50051"

	hoster := NewHoster(service, grpcAddr)
	hoster.HTTPAddr = httpAddr

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: time.Millisecond * 500,
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/hello?name=eleniums", httpAddr), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(httpReq)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	httpResp := pb.HelloResponse{}
	err = json.Unmarshal(body, &httpResp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.Equal(t, "Hello eleniums!", httpResp.Greeting)
}
