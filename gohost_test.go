package gohost

import (
	"testing"

	"github.com/eleniums/gohost/examples/hello"

	assert "github.com/stretchr/testify/require"
)

func Test_ServeGRPC_NilService(t *testing.T) {
	// arrange
	grpcAddr := "127.0.0.1:50051"

	// act
	err := ServeGRPC(nil, grpcAddr, nil)

	// assert
	assert.Error(t, err)
}

func Test_ServeGRPC_EmptyAddress(t *testing.T) {
	// arrange
	service := hello.NewService()

	// act
	err := ServeGRPC(service, "", nil)

	// assert
	assert.Error(t, err)
}
