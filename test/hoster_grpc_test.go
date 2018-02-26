package test

// import (
// 	"crypto/tls"
// 	"testing"
// 	"time"

// 	"github.com/eleniums/gohost/examples/test"
// 	"golang.org/x/net/context"
// 	g "google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"

// 	pb "github.com/eleniums/gohost/examples/test/proto"
// 	assert "github.com/stretchr/testify/require"
// )

// func Test_ServeGRPC_Successful(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	expectedValue := "test"

// 	// act - start the service
// 	go ServeGRPC(service, grpcAddr, nil)

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service
// 	conn, err := g.Dial(grpcAddr, g.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	req := pb.SendRequest{
// 		Value: expectedValue,
// 	}
// 	resp, err := client.Echo(context.Background(), &req)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, expectedValue, resp.Echo)
// }

// func Test_ServeGRPC_NilService(t *testing.T) {
// 	// arrange
// 	grpcAddr := getAddr(t)

// 	// act
// 	err := ServeGRPC(nil, grpcAddr, nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPC_EmptyGRPCAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()

// 	// act
// 	err := ServeGRPC(service, "", nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPC_FailListen(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := "badaddress"

// 	// act - start the service
// 	err := ServeGRPC(service, grpcAddr, nil)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPCWithTLS_Successful(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	expectedValue := "test"

// 	// act - start the service
// 	go ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service
// 	conn, err := g.Dial(grpcAddr, g.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	req := pb.SendRequest{
// 		Value: expectedValue,
// 	}
// 	resp, err := client.Echo(context.Background(), &req)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, resp)
// 	assert.Equal(t, expectedValue, resp.Echo)
// }

// func Test_ServeGRPCWithTLS_NilService(t *testing.T) {
// 	// arrange
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeGRPCWithTLS(nil, grpcAddr, nil, certFile, keyFile)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPCWithTLS_EmptyGRPCAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeGRPCWithTLS(service, "", nil, certFile, keyFile)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPCWithTLS_EmptyCertFile(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)
// 	keyFile := "./testdata/test.key"

// 	// act
// 	err := ServeGRPCWithTLS(service, grpcAddr, nil, "", keyFile)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPCWithTLS_EmptyKeyFile(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/test.crt"

// 	// act
// 	err := ServeGRPCWithTLS(service, grpcAddr, nil, certFile, "")

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPCWithTLS_BadCertFile(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)
// 	certFile := "./testdata/badcert.crt"
// 	keyFile := "./testdata/test.key"

// 	// act - start the service
// 	err := ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_ServeGRPCWithTLS_FailListen(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := "badaddress"
// 	certFile := "./testdata/test.crt"
// 	keyFile := "./testdata/test.key"

// 	// act - start the service
// 	err := ServeGRPCWithTLS(service, grpcAddr, nil, certFile, keyFile)

// 	// assert
// 	assert.Error(t, err)
// }
