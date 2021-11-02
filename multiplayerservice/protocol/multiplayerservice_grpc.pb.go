// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package nexus2

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MultiplayerServiceClient is the client API for MultiplayerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MultiplayerServiceClient interface {
	Login(ctx context.Context, opts ...grpc.CallOption) (MultiplayerService_LoginClient, error)
}

type multiplayerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMultiplayerServiceClient(cc grpc.ClientConnInterface) MultiplayerServiceClient {
	return &multiplayerServiceClient{cc}
}

func (c *multiplayerServiceClient) Login(ctx context.Context, opts ...grpc.CallOption) (MultiplayerService_LoginClient, error) {
	stream, err := c.cc.NewStream(ctx, &MultiplayerService_ServiceDesc.Streams[0], "/proto.MultiplayerService/Login", opts...)
	if err != nil {
		return nil, err
	}
	x := &multiplayerServiceLoginClient{stream}
	return x, nil
}

type MultiplayerService_LoginClient interface {
	Send(*ClientMessage) error
	Recv() (*ServerMessage, error)
	grpc.ClientStream
}

type multiplayerServiceLoginClient struct {
	grpc.ClientStream
}

func (x *multiplayerServiceLoginClient) Send(m *ClientMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *multiplayerServiceLoginClient) Recv() (*ServerMessage, error) {
	m := new(ServerMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MultiplayerServiceServer is the server API for MultiplayerService service.
// All implementations must embed UnimplementedMultiplayerServiceServer
// for forward compatibility
type MultiplayerServiceServer interface {
	Login(MultiplayerService_LoginServer) error
	mustEmbedUnimplementedMultiplayerServiceServer()
}

// UnimplementedMultiplayerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMultiplayerServiceServer struct {
}

func (UnimplementedMultiplayerServiceServer) Login(MultiplayerService_LoginServer) error {
	return status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedMultiplayerServiceServer) mustEmbedUnimplementedMultiplayerServiceServer() {}

// UnsafeMultiplayerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MultiplayerServiceServer will
// result in compilation errors.
type UnsafeMultiplayerServiceServer interface {
	mustEmbedUnimplementedMultiplayerServiceServer()
}

func RegisterMultiplayerServiceServer(s grpc.ServiceRegistrar, srv MultiplayerServiceServer) {
	s.RegisterService(&MultiplayerService_ServiceDesc, srv)
}

func _MultiplayerService_Login_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MultiplayerServiceServer).Login(&multiplayerServiceLoginServer{stream})
}

type MultiplayerService_LoginServer interface {
	Send(*ServerMessage) error
	Recv() (*ClientMessage, error)
	grpc.ServerStream
}

type multiplayerServiceLoginServer struct {
	grpc.ServerStream
}

func (x *multiplayerServiceLoginServer) Send(m *ServerMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *multiplayerServiceLoginServer) Recv() (*ClientMessage, error) {
	m := new(ClientMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MultiplayerService_ServiceDesc is the grpc.ServiceDesc for MultiplayerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MultiplayerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MultiplayerService",
	HandlerType: (*MultiplayerServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Login",
			Handler:       _MultiplayerService_Login_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "multiplayerservice/protocol/multiplayerservice.proto",
}
