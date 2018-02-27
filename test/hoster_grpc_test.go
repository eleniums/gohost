package test

import (
	"crypto/tls"
	"math"
	"testing"
	"time"

	"github.com/eleniums/gohost"
	"github.com/eleniums/gohost/examples/test"
	pb "github.com/eleniums/gohost/examples/test/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	assert "github.com/stretchr/testify/require"
)

func Test_Hoster_ListenAndServe_GRPC_Successful(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: expectedValue,
	}
	grpcResp, err := client.Echo(context.Background(), &grpcReq)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, expectedValue, grpcResp.Echo)
}

func Test_Hoster_ListenAndServe_GRPC_WithTLS(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.CertFile = "./testdata/test.crt"
	hoster.KeyFile = "./testdata/test.key"

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: expectedValue,
	}
	grpcResp, err := client.Echo(context.Background(), &grpcReq)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, expectedValue, grpcResp.Echo)
}

func Test_Hoster_ListenAndServe_GRPC_EmptyAddress(t *testing.T) {
	// arrange
	service := test.NewService()

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = ""
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	// act - start the service
	err := hoster.ListenAndServe()

	// assert
	assert.Error(t, err)
}

func Test_Hoster_ListenAndServe_GRPC_MaxRecvMsgSize_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	largeValue := string(make([]byte, largeMessageLength))

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.MaxRecvMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: largeValue,
	}
	grpcResp, err := client.Send(context.Background(), &grpcReq, grpc.MaxCallSendMsgSize(math.MaxInt32))

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.True(t, grpcResp.Success)
}

func Test_Hoster_ListenAndServe_GRPC_MaxRecvMsgSize_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	largeValue := string(make([]byte, largeMessageLength))

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.MaxRecvMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: largeValue,
	}
	grpcResp, err := client.Send(context.Background(), &grpcReq, grpc.MaxCallSendMsgSize(math.MaxInt32))

	// assert
	assert.Error(t, err)
	assert.Nil(t, grpcResp)
}

func Test_Hoster_ListenAndServe_GRPC_MaxSendMsgSize_Pass(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.MaxSendMsgSize = math.MaxInt32

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.LargeRequest{
		Length: largeMessageLength,
	}
	grpcResp, err := client.Large(context.Background(), &grpcReq, grpc.MaxCallRecvMsgSize(math.MaxInt32))

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, largeMessageLength, len(grpcResp.Echo))
}

func Test_Hoster_ListenAndServe_GRPC_MaxSendMsgSize_Fail(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	hoster.MaxSendMsgSize = 1

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.LargeRequest{
		Length: largeMessageLength,
	}
	grpcResp, err := client.Large(context.Background(), &grpcReq, grpc.MaxCallRecvMsgSize(math.MaxInt32))

	// assert
	assert.Error(t, err)
	assert.Nil(t, grpcResp)
}

func Test_Hoster_ListenAndServe_GRPC_UnaryInterceptor(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	count := 1
	hoster.UnaryInterceptors = append(hoster.UnaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		assert.Equal(t, 1, count)
		count++
		return handler(ctx, req)
	})
	hoster.UnaryInterceptors = append(hoster.UnaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		assert.Equal(t, 2, count)
		count++
		return handler(ctx, req)
	})
	hoster.UnaryInterceptors = append(hoster.UnaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		assert.Equal(t, 3, count)
		count++
		return handler(ctx, req)
	})

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcReq := pb.SendRequest{
		Value: expectedValue,
	}
	grpcResp, err := client.Echo(context.Background(), &grpcReq)

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, grpcResp)
	assert.Equal(t, expectedValue, grpcResp.Echo)
}

func Test_Hoster_ListenAndServe_GRPC_StreamInterceptor(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

	count := 1
	hoster.StreamInterceptors = append(hoster.StreamInterceptors, func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		assert.Equal(t, 1, count)
		count++
		return handler(srv, stream)
	})
	hoster.StreamInterceptors = append(hoster.StreamInterceptors, func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		assert.Equal(t, 2, count)
		count++
		return handler(srv, stream)
	})
	hoster.StreamInterceptors = append(hoster.StreamInterceptors, func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		assert.Equal(t, 3, count)
		count++
		return handler(srv, stream)
	})

	// act - start the service
	go hoster.ListenAndServe()

	// make sure service has time to start
	time.Sleep(serviceStartDelay)

	// call the service at the gRPC endpoint
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	client := pb.NewTestServiceClient(conn)
	grpcResp, err := client.Stream(context.Background())
	assert.NoError(t, err)

	// send some values
	err = grpcResp.Send(&pb.SendRequest{
		Value: "value1",
	})
	assert.NoError(t, err)

	err = grpcResp.Send(&pb.SendRequest{
		Value: "value2",
	})
	assert.NoError(t, err)

	err = grpcResp.Send(&pb.SendRequest{
		Value: "value3",
	})
	assert.NoError(t, err)

	// close out the stream
	testResp, err := grpcResp.CloseAndRecv()

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, testResp)
	assert.True(t, testResp.Success)
}

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
