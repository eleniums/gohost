// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/test.proto

/*
Package test is a generated protocol buffer package.

It is generated from these files:
	proto/test.proto

It has these top-level messages:
	SendRequest
	LargeRequest
	TestResponse
	EchoResponse
*/
package test

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Send request.
type SendRequest struct {
	// Value to send.
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *SendRequest) Reset()                    { *m = SendRequest{} }
func (m *SendRequest) String() string            { return proto.CompactTextString(m) }
func (*SendRequest) ProtoMessage()               {}
func (*SendRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SendRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

// Large request.
type LargeRequest struct {
	// Length of string to return in response.
	Length int64 `protobuf:"varint,1,opt,name=length" json:"length,omitempty"`
}

func (m *LargeRequest) Reset()                    { *m = LargeRequest{} }
func (m *LargeRequest) String() string            { return proto.CompactTextString(m) }
func (*LargeRequest) ProtoMessage()               {}
func (*LargeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *LargeRequest) GetLength() int64 {
	if m != nil {
		return m.Length
	}
	return 0
}

// Test response.
type TestResponse struct {
	// True if operation was a success.
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *TestResponse) Reset()                    { *m = TestResponse{} }
func (m *TestResponse) String() string            { return proto.CompactTextString(m) }
func (*TestResponse) ProtoMessage()               {}
func (*TestResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TestResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

// Echo response.
type EchoResponse struct {
	// Echo from service.
	Echo string `protobuf:"bytes,1,opt,name=echo" json:"echo,omitempty"`
}

func (m *EchoResponse) Reset()                    { *m = EchoResponse{} }
func (m *EchoResponse) String() string            { return proto.CompactTextString(m) }
func (*EchoResponse) ProtoMessage()               {}
func (*EchoResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *EchoResponse) GetEcho() string {
	if m != nil {
		return m.Echo
	}
	return ""
}

func init() {
	proto.RegisterType((*SendRequest)(nil), "test.SendRequest")
	proto.RegisterType((*LargeRequest)(nil), "test.LargeRequest")
	proto.RegisterType((*TestResponse)(nil), "test.TestResponse")
	proto.RegisterType((*EchoResponse)(nil), "test.EchoResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for TestService service

type TestServiceClient interface {
	// Echo the value in the request back in the response.
	Echo(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	// Send the value in the request.
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*TestResponse, error)
	// Large will return a large response message.
	Large(ctx context.Context, in *LargeRequest, opts ...grpc.CallOption) (*EchoResponse, error)
}

type testServiceClient struct {
	cc *grpc.ClientConn
}

func NewTestServiceClient(cc *grpc.ClientConn) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) Echo(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := grpc.Invoke(ctx, "/test.TestService/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*TestResponse, error) {
	out := new(TestResponse)
	err := grpc.Invoke(ctx, "/test.TestService/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) Large(ctx context.Context, in *LargeRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := grpc.Invoke(ctx, "/test.TestService/Large", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for TestService service

type TestServiceServer interface {
	// Echo the value in the request back in the response.
	Echo(context.Context, *SendRequest) (*EchoResponse, error)
	// Send the value in the request.
	Send(context.Context, *SendRequest) (*TestResponse, error)
	// Large will return a large response message.
	Large(context.Context, *LargeRequest) (*EchoResponse, error)
}

func RegisterTestServiceServer(s *grpc.Server, srv TestServiceServer) {
	s.RegisterService(&_TestService_serviceDesc, srv)
}

func _TestService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.TestService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Echo(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.TestService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_Large_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LargeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).Large(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.TestService/Large",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).Large(ctx, req.(*LargeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "test.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _TestService_Echo_Handler,
		},
		{
			MethodName: "Send",
			Handler:    _TestService_Send_Handler,
		},
		{
			MethodName: "Large",
			Handler:    _TestService_Large_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/test.proto",
}

func init() { proto.RegisterFile("proto/test.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x41, 0x4a, 0xc3, 0x40,
	0x14, 0x86, 0x89, 0xa6, 0xb5, 0x7d, 0xc9, 0xa2, 0x7d, 0x88, 0x94, 0xe2, 0x42, 0x46, 0x90, 0xae,
	0x1a, 0xd4, 0x03, 0x08, 0x82, 0x3b, 0x57, 0x53, 0x2f, 0x30, 0x8e, 0x8f, 0x24, 0x10, 0x66, 0x62,
	0xde, 0x24, 0x07, 0xf0, 0x0a, 0x1e, 0xad, 0x57, 0xf0, 0x20, 0x32, 0x93, 0x44, 0xb2, 0x90, 0xee,
	0xde, 0x0b, 0xdf, 0xfb, 0xf2, 0xcf, 0x0f, 0xab, 0xba, 0xb1, 0xce, 0x66, 0x8e, 0xd8, 0xed, 0xc3,
	0x88, 0xb1, 0x9f, 0xb7, 0xd7, 0xb9, 0xb5, 0x79, 0x45, 0x99, 0xaa, 0xcb, 0x4c, 0x19, 0x63, 0x9d,
	0x72, 0xa5, 0x35, 0xdc, 0x33, 0xe2, 0x16, 0x92, 0x03, 0x99, 0x0f, 0x49, 0x9f, 0x2d, 0xb1, 0xc3,
	0x4b, 0x98, 0x75, 0xaa, 0x6a, 0x69, 0x13, 0xdd, 0x44, 0xbb, 0xa5, 0xec, 0x17, 0x71, 0x07, 0xe9,
	0xab, 0x6a, 0x72, 0x1a, 0xa9, 0x2b, 0x98, 0x57, 0x64, 0x72, 0x57, 0x04, 0xec, 0x5c, 0x0e, 0x9b,
	0xd8, 0x41, 0xfa, 0x46, 0xec, 0x24, 0x71, 0x6d, 0x0d, 0x13, 0x6e, 0xe0, 0x82, 0x5b, 0xad, 0x89,
	0x39, 0x80, 0x0b, 0x39, 0xae, 0x42, 0x40, 0xfa, 0xa2, 0x0b, 0xfb, 0x47, 0x22, 0xc4, 0xa4, 0x0b,
	0x3b, 0xfc, 0x36, 0xcc, 0x0f, 0xc7, 0x08, 0x12, 0xaf, 0x3b, 0x50, 0xd3, 0x95, 0x9a, 0xf0, 0x09,
	0x62, 0x7f, 0x83, 0xeb, 0x7d, 0x78, 0xe3, 0x24, 0xf6, 0x16, 0xfb, 0x4f, 0x53, 0xa5, 0x58, 0x7d,
	0x1d, 0x7f, 0xbe, 0xcf, 0x00, 0x17, 0x59, 0x77, 0x9f, 0x79, 0xa1, 0x17, 0xf8, 0xa3, 0x13, 0x82,
	0x69, 0xfa, 0x51, 0x20, 0x82, 0x80, 0xfd, 0xe1, 0x33, 0xcc, 0x42, 0x0f, 0x38, 0xe0, 0xd3, 0x52,
	0xfe, 0xcd, 0xb0, 0x0e, 0x8a, 0x04, 0x97, 0x5e, 0x51, 0x79, 0xfa, 0x7d, 0x1e, 0x7a, 0x7f, 0xfc,
	0x0d, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xc7, 0x93, 0x68, 0xaf, 0x01, 0x00, 0x00,
}
