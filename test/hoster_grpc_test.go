package test

// import (
// 	"testing"

// 	"github.com/eleniums/gohost/examples/test"

// 	assert "github.com/stretchr/testify/require"
// )

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
