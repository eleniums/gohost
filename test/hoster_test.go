package test

import (
	"flag"
	"net"
	"os"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

const (
	largeMessageLength = 1000
)

var (
	serviceStartDelay = time.Millisecond * 100
	httpClientTimeout = time.Millisecond * 5000
)

func TestMain(m *testing.M) {
	flag.DurationVar(&serviceStartDelay, "service-start-delay", serviceStartDelay, "time to delay in milliseconds so test service can start")
	flag.DurationVar(&httpClientTimeout, "http-client-timeout", httpClientTimeout, "http client timeout in milliseconds")
	flag.Parse()

	os.Exit(m.Run())
}

// getAddr is a helper function that will retrieve a 127.0.0.1 address with an open port.
func getAddr(t *testing.T) string {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer lis.Close()

	return lis.Addr().String()
}
