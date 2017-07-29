package gohost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/eleniums/gohost/examples/test"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/test/proto"
	assert "github.com/stretchr/testify/require"
)

const (
	largeMessageLength = 10000000
)

func Test_Hoster_ListenAndServe_GRPCEndpoint(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50051"

	expectedValue := "test"

	hoster := NewHoster(service, grpcAddr)

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: expectedValue,
	}
	grpcResp, err := client.Echo(context.Background(), &grpcReq)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, expectedValue, grpcResp.Echo)
}

func Test_Hoster_ListenAndServe_HTTPEndpoint(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50052"

	expectedValue := "test"

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
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/echo?value="+expectedValue, httpAddr), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(httpReq)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	httpResp := pb.EchoResponse{}
	err = json.Unmarshal(body, &httpResp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.Equal(t, expectedValue, httpResp.Echo)
}

func Test_Hoster_ListenAndServe_MaxRecvMsgSize_GRPC_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50053"

	largeValue := string(make([]byte, largeMessageLength))

	hoster := NewHoster(service, grpcAddr)
	hoster.MaxRecvMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: largeValue,
	}
	grpcResp, err := client.Send(context.Background(), &grpcReq, grpc.MaxCallSendMsgSize(math.MaxInt32))

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.True(t, grpcResp.Success)
}

func Test_Hoster_ListenAndServe_MaxRecvMsgSize_GRPC_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50054"

	largeValue := string(make([]byte, largeMessageLength))

	hoster := NewHoster(service, grpcAddr)
	hoster.MaxRecvMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: largeValue,
	}
	grpcResp, err := client.Send(context.Background(), &grpcReq, grpc.MaxCallSendMsgSize(math.MaxInt32))

	// assert
	assert.Error(t, err)
	assert.Nil(t, grpcResp)
}

func Test_Hoster_ListenAndServe_MaxRecvMsgSize_HTTP_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9090"
	grpcAddr := "127.0.0.1:50055"

	largeValue := string(make([]byte, largeMessageLength))

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
	httpReq := pb.SendRequest{
		Value: largeValue,
	}
	payload, err := json.Marshal(&httpReq)
	assert.NoError(t, err)
	postReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%v/v1/send", httpAddr), bytes.NewBuffer(payload))
	assert.NoError(t, err)
	doResp, err := httpClient.Do(postReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, doResp.StatusCode)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	httpResp := pb.TestResponse{}
	err = json.Unmarshal(body, &httpResp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.True(t, httpResp.Success)
}

func Test_Hoster_ListenAndServe_MaxSendMsgSize_GRPC_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50056"

	hoster := NewHoster(service, grpcAddr)
	hoster.MaxSendMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.LargeRequest{
		Length: largeMessageLength,
	}
	grpcResp, err := client.Large(context.Background(), &grpcReq, grpc.MaxCallRecvMsgSize(math.MaxInt32))

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, largeMessageLength, len(grpcResp.Echo))
}

func Test_Hoster_ListenAndServe_MaxSendMsgSize_GRPC_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := "127.0.0.1:50057"

	hoster := NewHoster(service, grpcAddr)
	hoster.MaxSendMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.LargeRequest{
		Length: largeMessageLength,
	}
	grpcResp, err := client.Large(context.Background(), &grpcReq, grpc.MaxCallRecvMsgSize(math.MaxInt32))

	// assert
	assert.Error(t, err)
	assert.Nil(t, grpcResp)
}

func Test_Hoster_ListenAndServe_MaxSendMsgSize_HTTP_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9091"
	grpcAddr := "127.0.0.1:50058"

	hoster := NewHoster(service, grpcAddr)
	hoster.HTTPAddr = httpAddr
	hoster.MaxSendMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: time.Millisecond * 1000,
	}
	postReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/large?length=%v", httpAddr, largeMessageLength), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(postReq)
	assert.NoError(t, err)
	assert.Equal(t, 200, doResp.StatusCode)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	httpResp := pb.EchoResponse{}
	err = json.Unmarshal(body, &httpResp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.Equal(t, largeMessageLength, len(httpResp.Echo))
}

func Test_Hoster_ListenAndServe_MaxSendMsgSize_HTTP_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "127.0.0.1:9092"
	grpcAddr := "127.0.0.1:50059"

	hoster := NewHoster(service, grpcAddr)
	hoster.HTTPAddr = httpAddr
	hoster.MaxSendMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: time.Millisecond * 1000,
	}
	postReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/large?length=%v", httpAddr, largeMessageLength), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(postReq)
	assert.NoError(t, err)
	assert.NotEqual(t, 200, doResp.StatusCode)
}
