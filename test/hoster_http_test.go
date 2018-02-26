package test

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	h "net/http"
	"testing"
	"time"

	"github.com/eleniums/gohost/examples/test"

	pb "github.com/eleniums/gohost/examples/test/proto"
	assert "github.com/stretchr/testify/require"
)

func Test_ServeHTTP_Successful(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	expectedValue := "test"

	// act - start the service
	go ServeGRPC(service, grpcAddr, nil)
	go ServeHTTP(service, httpAddr, grpcAddr, false, nil)

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service
	httpClient := h.Client{
		Timeout: httpClientTimeout,
	}
	req, err := h.NewRequest(h.MethodGet, fmt.Sprintf("http://%v/v1/echo?value="+expectedValue, httpAddr), nil)
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

func Test_ServeHTTP_EnableCORS(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	expectedValue := "test"

	// act - start the service
	go ServeGRPC(service, grpcAddr, nil)
	go ServeHTTP(service, httpAddr, grpcAddr, true, nil)

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service
	httpClient := h.Client{
		Timeout: httpClientTimeout,
	}
	req, err := h.NewRequest(h.MethodGet, fmt.Sprintf("http://%v/v1/echo?value="+expectedValue, httpAddr), nil)
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
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)

	// act
	err := ServeHTTP(nil, httpAddr, grpcAddr, false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTP_EmptyHTTPAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	// act
	err := ServeHTTP(service, "", grpcAddr, false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTP_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)

	// act
	err := ServeHTTP(service, httpAddr, "", false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTP_FailListen(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "badaddress"
	grpcAddr := getAddr(t)

	// act - start the service
	go ServeGRPC(service, grpcAddr, nil)
	err := ServeHTTP(service, httpAddr, grpcAddr, false, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_Successful(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)
	certFile := "./testdata/test.crt"
	keyFile := "./testdata/test.key"

	expectedValue := "test"

	// act - start the service
	go ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)
	go ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, certFile, keyFile, true)

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service
	httpClient := h.Client{
		Timeout: httpClientTimeout,
		Transport: &h.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := h.NewRequest(h.MethodGet, fmt.Sprintf("https://%v/v1/echo?value="+expectedValue, httpAddr), nil)
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

func Test_ServeHTTPWithTLS_EnableCORS(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)
	certFile := "./testdata/test.crt"
	keyFile := "./testdata/test.key"

	expectedValue := "test"

	// act - start the service
	go ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)
	go ServeHTTPWithTLS(service, httpAddr, grpcAddr, true, nil, certFile, keyFile, true)

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service
	httpClient := h.Client{
		Timeout: httpClientTimeout,
		Transport: &h.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := h.NewRequest(h.MethodGet, fmt.Sprintf("https://%v/v1/echo?value="+expectedValue, httpAddr), nil)
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

func Test_ServeHTTPWithTLS_NilService(t *testing.T) {
	// arrange
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)
	certFile := "./testdata/test.crt"
	keyFile := "./testdata/test.key"

	// act
	err := ServeHTTPWithTLS(nil, httpAddr, grpcAddr, false, nil, certFile, keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyHTTPAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)
	certFile := "./testdata/test.crt"
	keyFile := "./testdata/test.key"

	// act
	err := ServeHTTPWithTLS(service, "", grpcAddr, false, nil, certFile, keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyGRPCAddress(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	certFile := "./testdata/test.crt"
	keyFile := "./testdata/test.key"

	// act
	err := ServeHTTPWithTLS(service, httpAddr, "", false, nil, certFile, keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyCertFile(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)
	keyFile := "./testdata/test.key"

	// act
	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, "", keyFile, false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_EmptyKeyFile(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := getAddr(t)
	grpcAddr := getAddr(t)
	certFile := "./testdata/test.crt"

	// act
	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, certFile, "", false)

	// assert
	assert.Error(t, err)
}

func Test_ServeHTTPWithTLS_FailListen(t *testing.T) {
	// arrange
	service := test.NewService()
	httpAddr := "badaddress"
	grpcAddr := getAddr(t)
	certFile := "./testdata/test.crt"
	keyFile := "./testdata/test.key"

	// act - start the service
	go ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)
	err := ServeHTTPWithTLS(service, httpAddr, grpcAddr, false, nil, certFile, keyFile, true)

	// assert
	assert.Error(t, err)
}
