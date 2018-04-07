package test

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

	"github.com/eleniums/gohost"
	"github.com/eleniums/gohost/examples/test"
	"google.golang.org/grpc"

	pb "github.com/eleniums/gohost/examples/test/proto"
	assert "github.com/stretchr/testify/require"
)

func Test_Hoster_ListenAndServe_HTTP_Successful(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/echo?value=%v", httpAddr, expectedValue), nil)
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

func Test_Hoster_ListenAndServe_HTTP_WithTLS(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.CertFile = "../testdata/test.crt"
	hoster.KeyFile = "../testdata/test.key"
	hoster.InsecureSkipVerify = true

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%v/v1/echo?value=%v", httpAddr, expectedValue), nil)
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

func Test_Hoster_ListenAndServe_HTTP_EmptyAddress(t *testing.T) {
	// arrange
	hoster := gohost.NewHoster()

	hoster.HTTPAddr = ""
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_HTTP_InvalidAddress(t *testing.T) {
	// arrange
	hoster := gohost.NewHoster()

	hoster.HTTPAddr = "badaddress"
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_HTTP_InvalidCertFile(t *testing.T) {
	// arrange
	hoster := gohost.NewHoster()
	httpAddr := getAddr(t)

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.CertFile = "../testdata/badcert.crt"
	hoster.KeyFile = "../testdata/test.key"

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_HTTP_InvalidKeyFile(t *testing.T) {
	// arrange
	hoster := gohost.NewHoster()
	httpAddr := getAddr(t)

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.CertFile = "../testdata/test.crt"
	hoster.KeyFile = "../testdata/badkey.key"

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_HTTP_MaxRecvMsgSize_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	largeValue := string(make([]byte, largeMessageLength))

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.MaxRecvMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
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

func Test_Hoster_ListenAndServe_HTTP_MaxRecvMsgSize_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	largeValue := string(make([]byte, largeMessageLength))

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.MaxRecvMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
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
	assert.NotEqual(t, 200, doResp.StatusCode)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	httpResp := pb.TestResponse{}
	err = json.Unmarshal(body, &httpResp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.False(t, httpResp.Success)
}

func Test_Hoster_ListenAndServe_HTTP_MaxSendMsgSize_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.MaxSendMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
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

func Test_Hoster_ListenAndServe_HTTP_MaxSendMsgSize_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.HTTPAddr = httpAddr
	hoster.RegisterHTTPEndpoint(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.MaxSendMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
	}
	postReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/large?length=%v", httpAddr, largeMessageLength), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(postReq)
	assert.NoError(t, err)
	assert.NotEqual(t, 200, doResp.StatusCode)
	body, err := ioutil.ReadAll(doResp.Body)
	assert.NoError(t, err)
	httpResp := pb.EchoResponse{}
	err = json.Unmarshal(body, &httpResp)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, httpResp)
	assert.Empty(t, httpResp.Echo)
}
