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

func Test_Hoster_ListenAndServe_DebugEndpoint_Pprof(t *testing.T) {
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

func Test_Hoster_ListenAndServe_DebugEndpoint_Vars(t *testing.T) {
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
