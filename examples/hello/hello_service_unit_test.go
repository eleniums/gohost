package hello

import (
	"testing"

	"golang.org/x/net/context"

	pb "github.com/eleniums/gohost/examples/hello/proto"
	assert "github.com/stretchr/testify/require"
)

func Test_Service_Hello_WithName(t *testing.T) {
	// arrange
	service := NewService()
	ctx := context.TODO()

	req := pb.HelloRequest{
		Name: "eleniums",
	}

	// act
	resp, err := service.Hello(ctx, &req)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Greeting)
	assert.Equal(t, "Hello eleniums!", resp.Greeting)
}

func Test_Service_Hello_NoName(t *testing.T) {
	// arrange
	service := NewService()
	ctx := context.TODO()

	req := pb.HelloRequest{}

	// act
	resp, err := service.Hello(ctx, &req)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Greeting)
	assert.Equal(t, "Hello!", resp.Greeting)
}
