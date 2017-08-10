package gohost

import (
	"bytes"
	"crypto/tls"
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
	"google.golang.org/grpc/credentials"

	pb "github.com/eleniums/gohost/examples/test/proto"
	assert "github.com/stretchr/testify/require"
)

const (
	largeMessageLength = 1000
)

func Test_Hoster_ListenAndServe_GRPCEndpoint(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

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

func Test_Hoster_ListenAndServe_GRPCEndpoint_WithTLS(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := NewHoster(service, grpcAddr)
	hoster.CertFile = "./testdata/test.crt"
	hoster.KeyFile = "./testdata/test.key"

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
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
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

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

func Test_Hoster_ListenAndServe_HTTPEndpoint_WithTLS(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := NewHoster(service, grpcAddr)
	hoster.HTTPAddr = httpAddr
	hoster.CertFile = "./testdata/test.crt"
	hoster.KeyFile = "./testdata/test.key"
	hoster.InsecureSkipVerify = true

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(time.Millisecond * 100)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: time.Millisecond * 500,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%v/v1/echo?value="+expectedValue, httpAddr), nil)
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

func Test_Hoster_ListenAndServe_Logger(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := NewHoster(service, grpcAddr)

	loggedValue := ""
	hoster.Logger = func(format string, v ...interface{}) {
		loggedValue = fmt.Sprintf(format, v...)
	}

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
	_, err = client.Echo(context.Background(), &grpcReq)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, loggedValue)
}

func Test_Hoster_ListenAndServe_NilService(t *testing.T) {
	// arrange
	grpcAddr := getAddr(t)

	hoster := NewHoster(nil, grpcAddr)

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()

	hoster := NewHoster(service, "")

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_DoesNotImplementHTTPInterface(t *testing.T) {
	// arrange
	service := test.NewGRPCService()
	grpcAddr := getAddr(t)
	httpAddr := getAddr(t)

	hoster := NewHoster(service, grpcAddr)
	hoster.HTTPAddr = httpAddr

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_MaxRecvMsgSize_GRPC_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

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
	grpcAddr := getAddr(t)

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
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

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
	grpcAddr := getAddr(t)

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
	grpcAddr := getAddr(t)

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
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

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
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

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
