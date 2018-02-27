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

func Test_Hoster_ListenAndServe_HTTPEndpoint(t *testing.T) {
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
	hoster.RegisterHTTPHandler(pb.RegisterTestServiceHandlerFromEndpoint)

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

func Test_Hoster_ListenAndServe_HTTPEndpoint_WithTLS(t *testing.T) {
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
	hoster.RegisterHTTPHandler(pb.RegisterTestServiceHandlerFromEndpoint)

	hoster.CertFile = "./testdata/test.crt"
	hoster.KeyFile = "./testdata/test.key"
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

func Test_Hoster_ListenAndServe_MaxRecvMsgSize_HTTP_Pass(t *testing.T) {
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
	hoster.RegisterHTTPHandler(pb.RegisterTestServiceHandlerFromEndpoint)

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

func Test_Hoster_ListenAndServe_MaxRecvMsgSize_HTTP_Fail(t *testing.T) {
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
	hoster.RegisterHTTPHandler(pb.RegisterTestServiceHandlerFromEndpoint)

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

func Test_Hoster_ListenAndServe_MaxSendMsgSize_HTTP_Pass(t *testing.T) {
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
	hoster.RegisterHTTPHandler(pb.RegisterTestServiceHandlerFromEndpoint)

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

func Test_Hoster_ListenAndServe_MaxSendMsgSize_HTTP_Fail(t *testing.T) {
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
	hoster.RegisterHTTPHandler(pb.RegisterTestServiceHandlerFromEndpoint)

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

// func Test_ServeHTTP_EnableCORS(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	expectedValue := "test"

// 	// act - start the service
// 	go ServeGRPC(service, grpcAddr, nil)
// 	go ServeHTTP(service, httpAddr, grpcAddr, true, nil)

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service
// 	httpClient := h.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	req, err := h.NewRequest(h.MethodGet, fmt.Sprintf("http://%v/v1/echo?value="+expectedValue, httpAddr), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, doResp.StatusCode)
// 	body, err := ioutil.ReadAll(doResp.Body)
// 	assert.NoError(t, err)
// 	resp := pb.EchoResponse{}
// 	err = json.Unmarshal(body, &resp)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, expectedValue, resp.Echo)
// }

// func Test_ServeHTTP_NilService(t *testing.T) {
// 	// arrange
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	// act
// 	err := ServeHTTP(nil, httpAddr, grpcAddr, false, nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTP_EmptyHTTPAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	// act
// 	err := ServeHTTP(service, "", grpcAddr, false, nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTP_EmptyGRPCAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)

// 	// act
// 	err := ServeHTTP(service, httpAddr, "", false, nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTP_FailListen(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := "badaddress"
// 	grpcAddr := getAddr(t)

// 	// act - start the service
// 	go ServeGRPC(service, grpcAddr, nil)
// 	err := ServeHTTP(service, httpAddr, grpcAddr, false, nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTPWithTLS_EnableCORS(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	expectedValue := "test"

// 	// act - start the service
// 	go ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)
// 	go ServeHTTPWithTLS(service, httpAddr, grpcAddr, true, nil, certFile, keyFile, true)

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service
// 	httpClient := h.Client{
// 		Timeout: httpClientTimeout,
// 		Transport: &h.Transport{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: true,
// 			},
// 		},
// 	}
// 	req, err := h.NewRequest(h.MethodGet, fmt.Sprintf("https://%v/v1/echo?value="+expectedValue, httpAddr), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, doResp.StatusCode)
// 	body, err := ioutil.ReadAll(doResp.Body)
// 	assert.NoError(t, err)
// 	resp := pb.EchoResponse{}
// 	err = json.Unmarshal(body, &resp)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, expectedValue, resp.Echo)
// }

// func Test_ServeHTTPWithTLS_NilService(t *testing.T) {
// 	// arrange
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeHTTPWithTLS(nil, httpAddr, grpcAddr, false, nil, certFile, keyFile, false)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTPWithTLS_EmptyHTTPAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeHTTPWithTLS(service, "", grpcAddr, false, nil, certFile, keyFile, false)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTPWithTLS_EmptyGRPCAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeHTTPWithTLS(service, httpAddr, "", false, nil, certFile, keyFile, false)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTPWithTLS_EmptyCertFile(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, "", keyFile, false)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTPWithTLS_EmptyKeyFile(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"

// 	// act
// 	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, certFile, "", false)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeHTTPWithTLS_FailListen(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := "badaddress"
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act - start the service
// 	go ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)
// 	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, certFile, keyFile, true)

// 	// assert
// 	assert.Error(t, err)
// }
