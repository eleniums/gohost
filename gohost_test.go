package gohost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/eleniums/gohost/examples/test"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/test/proto"
	assert "github.com/stretchr/testify/require"
)

func Test_ServeGRPC_Successful(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50051"

	expectedValue := "test"

	// act - start the service
	go ServeGRPC(service, grpcAddr, nil)

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	req := pb.SendRequest{
		Value: expectedValue,
	}
	resp, err := client.Echo(context.Background(), &req)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedValue, resp.Echo)
}

func Test_ServeGRPC_NilService(t *testing.T) {
	// arrange
	grpcAddr := "127.0.0.1:50051"

	// act
	err := ServeGRPC(nil, grpcAddr, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeGRPC_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()

	// act
	err := ServeGRPC(service, "", nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeGRPCWithTLS_NilService(t *testing.T) {
	// arrange
	grpcAddr := "127.0.0.1:50051"
	certFile := "certfile"
	keyFile := "keyfile"

	// act
	err := ServeGRPCWithTLS(nil, grpcAddr, nil, certFile, keyFile)

	// assert
	assert.Error(t, err)
}

func Test_ServeGRPCWithTLS_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	certFile := "certfile"
	keyFile := "keyfile"

	// act
	err := ServeGRPCWithTLS(service, "", nil, certFile, keyFile)

	// assert
	assert.Error(t, err)
}

func Test_ServeGRPCWithTLS_EmptyCertFile(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50051"
	keyFile := "keyfile"

	// act
	err := ServeGRPCWithTLS(service, grpcAddr, nil, "", keyFile)

	// assert
	assert.Error(t, err)
}

func Test_ServeGRPCWithTLS_EmptyKeyFile(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50051"
	certFile := "certfile"

	// act
	err := ServeGRPCWithTLS(service, grpcAddr, nil, certFile, "")

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTP_Successful(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50051"

	expectedValue := "test"

	// act - start the service
	go ServeGRPC(service, grpcAddr, nil)
	go ServeHTTP(service, httpAddr, grpcAddr, false, nil)

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service
	httpClient := http.Client{
		Timeout: time.Millisecond * 500,
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/echo?value="+expectedValue, httpAddr), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, doResp.StatusCode)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	resp := pb.EchoResponse{}
	err = json.Unmarshal(body, &resp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedValue, resp.Echo)
}

func Test_ServeHTTP_NilService(t *testing.T) {
	// arrange
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50051"

	// act
	err := ServeHTTP(nil, httpAddr, grpcAddr, false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTP_EmptyHTTPAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50051"

	// act
	err := ServeHTTP(service, "", grpcAddr, false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTP_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"

	// act
	err := ServeHTTP(service, httpAddr, "", false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_NilService(t *testing.T) {
	// arrange
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50051"
	certFile := "certfile"
	keyFile := "keyfile"

	// act
	err := ServeHTTPWithTLS(nil, httpAddr, grpcAddr, false, nil, certFile, keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyHTTPAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50051"
	certFile := "certfile"
	keyFile := "keyfile"

	// act
	err := ServeHTTPWithTLS(service, "", grpcAddr, false, nil, certFile, keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"
	certFile := "certfile"
	keyFile := "keyfile"

	// act
	err := ServeHTTPWithTLS(service, httpAddr, "", false, nil, certFile, keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyCertFile(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50051"
	keyFile := "keyfile"

	// act
	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, "", keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyKeyFile(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50051"
	certFile := "certfile"

	// act
	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, certFile, "", false)

	// assert
	assert.Error(t, err)
}
