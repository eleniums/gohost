package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/eleniums/gohost"

	assert "github.com/stretchr/testify/require"
)

func Test_Hoster_ListenAndServe_Debug_Pprof(t *testing.T) {
	// arrange
	debugAddr := getAddr(t)

	hoster := gohost.NewHoster()

	hoster.DebugAddr = debugAddr
	hoster.EnableDebug = true

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/debug/pprof", debugAddr), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(httpReq)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(doResp.Body)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}

func Test_Hoster_ListenAndServe_Debug_Vars(t *testing.T) {
	// arrange
	debugAddr := getAddr(t)

	hoster := gohost.NewHoster()

	hoster.DebugAddr = debugAddr
	hoster.EnableDebug = true

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/debug/vars", debugAddr), nil)
	assert.NoError(t, err)
	doResp, err := httpClient.Do(httpReq)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(doResp.Body)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}

func Test_Hoster_ListenAndServe_Debug_Disabled(t *testing.T) {
	// arrange
	debugAddr := getAddr(t)

	hoster := gohost.NewHoster()

	hoster.DebugAddr = debugAddr
	hoster.EnableDebug = false

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the HTTP endpoint
	httpClient := http.Client{
		Timeout: httpClientTimeout,
	}
	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/debug/pprof", debugAddr), nil)
	assert.NoError(t, err)
	_, err = httpClient.Do(httpReq)
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_Debug_EmptyAddress(t *testing.T) {
	// arrange
	hoster := gohost.NewHoster()

	hoster.DebugAddr = ""
	hoster.EnableDebug = true

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}
