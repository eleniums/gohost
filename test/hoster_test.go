package test

import (
	"crypto/tls"
	"flag"
	"net"
	"os"
	"testing"
	"time"

	"github.com/eleniums/gohost"
	"github.com/eleniums/gohost/examples/test"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/eleniums/gohost/examples/test/proto"
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

func Test_Hoster_ListenAndServe_GRPCEndpoint(t *testing.T) {
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

func Test_Hoster_ListenAndServe_GRPCEndpoint_WithTLS(t *testing.T) {
	// arrange
	service := test.NewService()
	grpcAddr := getAddr(t)

	expectedValue := "test"

	hoster := gohost.NewHoster()
	hoster.GRPCAddr = grpcAddr
	hoster.CertFile = "./testdata/test.crt"
	hoster.KeyFile = "./testdata/test.key"

	hoster.RegisterGRPCEndpoint(func(s *grpc.Server) {
		pb.RegisterTestServiceServer(s, service)
	})

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

// func Test_Hoster_ListenAndServe_HTTPEndpoint(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	expectedValue := "test"

// 	hoster := gohost.NewHoster()
// 	hoster.HTTPAddr = httpAddr

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/echo?value=%v", httpAddr, expectedValue), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(httpReq)
// 	assert.NoError(t, err)
// 	body, err := ioutil.ReadAll(doResp.Body)
// 	assert.NoError(t, err)
// 	httpResp := pb.EchoResponse{}
// 	err = json.Unmarshal(body, &httpResp)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, httpResp)
// 	assert.Equal(t, expectedValue, httpResp.Echo)
// }

// func Test_Hoster_ListenAndServe_HTTPEndpoint_WithTLS(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	expectedValue := "test"

// 	hoster := gohost.NewHoster()
// 	hoster.HTTPAddr = httpAddr
// 	hoster.CertFile = "./testdata/test.crt"
// 	hoster.KeyFile = "./testdata/test.key"
// 	hoster.InsecureSkipVerify = true

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: true,
// 			},
// 		},
// 	}
// 	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%v/v1/echo?value=%v", httpAddr, expectedValue), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(httpReq)
// 	assert.NoError(t, err)
// 	body, err := ioutil.ReadAll(doResp.Body)
// 	assert.NoError(t, err)
// 	httpResp := pb.EchoResponse{}
// 	err = json.Unmarshal(body, &httpResp)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, httpResp)
// 	assert.Equal(t, expectedValue, httpResp.Echo)
// }

// func Test_Hoster_ListenAndServe_DebugEndpoint_Pprof(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	debugAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.DebugAddr = debugAddr

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/debug/pprof", debugAddr), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(httpReq)
// 	assert.NoError(t, err)
// 	body, err := ioutil.ReadAll(doResp.Body)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, body)
// }

// func Test_Hoster_ListenAndServe_DebugEndpoint_Vars(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	debugAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.DebugAddr = debugAddr

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	httpReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/debug/vars", debugAddr), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(httpReq)
// 	assert.NoError(t, err)
// 	body, err := ioutil.ReadAll(doResp.Body)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, body)
// }

// func Test_Hoster_ListenAndServe_NilService(t *testing.T) {
// 	// arrange
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()

// 	// act - start the service
// 	err := hoster.ListenAndServe()

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_Hoster_ListenAndServe_EmptyGRPCAddress(t *testing.T) {
// 	// arrange
// 	service := test.NewService()

// 	hoster := gohost.NewHoster()

// 	// act - start the service
// 	err := hoster.ListenAndServe()

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_Hoster_ListenAndServe_DoesNotImplementHTTPInterface(t *testing.T) {
// 	// arrange
// 	service := test.NewGRPCService()
// 	grpcAddr := getAddr(t)
// 	httpAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.HTTPAddr = httpAddr

// 	// act - start the service
// 	err := hoster.ListenAndServe()

// 	// assert
// 	assert.Error(t, err)
// }

// func Test_Hoster_ListenAndServe_MaxRecvMsgSize_GRPC_Pass(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	largeValue := string(make([]byte, largeMessageLength))

// 	hoster := gohost.NewHoster()
// 	hoster.MaxRecvMsgSize = math.MaxInt32

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the gRPC endpoint
// 	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	grpcReq := pb.SendRequest{
// 		Value: largeValue,
// 	}
// 	grpcResp, err := client.Send(context.Background(), &grpcReq, grpc.MaxCallSendMsgSize(math.MaxInt32))

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, grpcResp)
// 	assert.True(t, grpcResp.Success)
// }

// func Test_Hoster_ListenAndServe_MaxRecvMsgSize_GRPC_Fail(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	largeValue := string(make([]byte, largeMessageLength))

// 	hoster := gohost.NewHoster()
// 	hoster.MaxRecvMsgSize = 1

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the gRPC endpoint
// 	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	grpcReq := pb.SendRequest{
// 		Value: largeValue,
// 	}
// 	grpcResp, err := client.Send(context.Background(), &grpcReq, grpc.MaxCallSendMsgSize(math.MaxInt32))

// 	// assert
// 	assert.Error(t, err)
// 	assert.Nil(t, grpcResp)
// }

// func Test_Hoster_ListenAndServe_MaxRecvMsgSize_HTTP_Pass(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	largeValue := string(make([]byte, largeMessageLength))

// 	hoster := gohost.NewHoster()
// 	hoster.HTTPAddr = httpAddr

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	httpReq := pb.SendRequest{
// 		Value: largeValue,
// 	}
// 	payload, err := json.Marshal(&httpReq)
// 	assert.NoError(t, err)
// 	postReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%v/v1/send", httpAddr), bytes.NewBuffer(payload))
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(postReq)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, doResp.StatusCode)
// 	body, err := ioutil.ReadAll(doResp.Body)
// 	assert.NoError(t, err)
// 	httpResp := pb.TestResponse{}
// 	err = json.Unmarshal(body, &httpResp)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, httpResp)
// 	assert.True(t, httpResp.Success)
// }

// func Test_Hoster_ListenAndServe_MaxSendMsgSize_GRPC_Pass(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.MaxSendMsgSize = math.MaxInt32

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the gRPC endpoint
// 	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	grpcReq := pb.LargeRequest{
// 		Length: largeMessageLength,
// 	}
// 	grpcResp, err := client.Large(context.Background(), &grpcReq, grpc.MaxCallRecvMsgSize(math.MaxInt32))

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, grpcResp)
// 	assert.Equal(t, largeMessageLength, len(grpcResp.Echo))
// }

// func Test_Hoster_ListenAndServe_MaxSendMsgSize_GRPC_Fail(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.MaxSendMsgSize = 1

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the gRPC endpoint
// 	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	grpcReq := pb.LargeRequest{
// 		Length: largeMessageLength,
// 	}
// 	grpcResp, err := client.Large(context.Background(), &grpcReq, grpc.MaxCallRecvMsgSize(math.MaxInt32))

// 	// assert
// 	assert.Error(t, err)
// 	assert.Nil(t, grpcResp)
// }

// func Test_Hoster_ListenAndServe_MaxSendMsgSize_HTTP_Pass(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.HTTPAddr = httpAddr
// 	hoster.MaxSendMsgSize = math.MaxInt32

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	postReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/large?length=%v", httpAddr, largeMessageLength), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(postReq)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, doResp.StatusCode)
// 	body, err := ioutil.ReadAll(doResp.Body)
// 	assert.NoError(t, err)
// 	httpResp := pb.EchoResponse{}
// 	err = json.Unmarshal(body, &httpResp)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, httpResp)
// 	assert.Equal(t, largeMessageLength, len(httpResp.Echo))
// }

// func Test_Hoster_ListenAndServe_MaxSendMsgSize_HTTP_Fail(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	httpAddr := getAddr(t)
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()
// 	hoster.HTTPAddr = httpAddr
// 	hoster.MaxSendMsgSize = 1

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the HTTP endpoint
// 	httpClient := http.Client{
// 		Timeout: httpClientTimeout,
// 	}
// 	postReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/v1/large?length=%v", httpAddr, largeMessageLength), nil)
// 	assert.NoError(t, err)
// 	doResp, err := httpClient.Do(postReq)
// 	assert.NoError(t, err)
// 	assert.NotEqual(t, 200, doResp.StatusCode)
// }

// func Test_Hoster_ListenAndServe_UnaryInterceptor(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	expectedValue := "test"

// 	hoster := gohost.NewHoster()

// 	count := 1
// 	hoster.UnaryInterceptors = append(hoster.UnaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		assert.Equal(t, 1, count)
// 		count++
// 		return handler(ctx, req)
// 	})
// 	hoster.UnaryInterceptors = append(hoster.UnaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		assert.Equal(t, 2, count)
// 		count++
// 		return handler(ctx, req)
// 	})
// 	hoster.UnaryInterceptors = append(hoster.UnaryInterceptors, func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
// 		assert.Equal(t, 3, count)
// 		count++
// 		return handler(ctx, req)
// 	})

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the gRPC endpoint
// 	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	grpcReq := pb.SendRequest{
// 		Value: expectedValue,
// 	}
// 	grpcResp, err := client.Echo(context.Background(), &grpcReq)

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, grpcResp)
// 	assert.Equal(t, expectedValue, grpcResp.Echo)
// }

// func Test_Hoster_ListenAndServe_StreamInterceptor(t *testing.T) {
// 	// arrange
// 	service := test.NewService()
// 	grpcAddr := getAddr(t)

// 	hoster := gohost.NewHoster()

// 	count := 1
// 	hoster.StreamInterceptors = append(hoster.StreamInterceptors, func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		assert.Equal(t, 1, count)
// 		count++
// 		return handler(srv, stream)
// 	})
// 	hoster.StreamInterceptors = append(hoster.StreamInterceptors, func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		assert.Equal(t, 2, count)
// 		count++
// 		return handler(srv, stream)
// 	})
// 	hoster.StreamInterceptors = append(hoster.StreamInterceptors, func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 		assert.Equal(t, 3, count)
// 		count++
// 		return handler(srv, stream)
// 	})

// 	// act - start the service
// 	go hoster.ListenAndServe()

// 	// make sure service has time to start
// 	time.Sleep(serviceStartDelay)

// 	// call the service at the gRPC endpoint
// 	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
// 	assert.NoError(t, err)
// 	client := pb.NewTestServiceClient(conn)
// 	grpcResp, err := client.Stream(context.Background())
// 	assert.NoError(t, err)

// 	// send some values
// 	err = grpcResp.Send(&pb.SendRequest{
// 		Value: "value1",
// 	})
// 	assert.NoError(t, err)

// 	err = grpcResp.Send(&pb.SendRequest{
// 		Value: "value2",
// 	})
// 	assert.NoError(t, err)

// 	err = grpcResp.Send(&pb.SendRequest{
// 		Value: "value3",
// 	})
// 	assert.NoError(t, err)

// 	// close out the stream
// 	testResp, err := grpcResp.CloseAndRecv()

// 	// assert
// 	assert.NoError(t, err)
// 	assert.NotNil(t, testResp)
// 	assert.True(t, testResp.Success)
// }

// getAddr is a helper function that will retrieve a 127.0.0.1 address with an open port.
func getAddr(t *testing.T) string {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer lis.Close()

	return lis.Addr().String()
}
