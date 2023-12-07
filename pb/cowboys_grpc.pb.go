// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: github.com/sumlookup/cowboys/pb/cowboys.proto

package cowboys

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

// CowboysServiceClient is the client API for CowboysService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CowboysServiceClient interface {
	Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error)
	ReloadDefaultCowboys(ctx context.Context, in *ReloadDefaultCowboysRequest, opts ...grpc.CallOption) (*ReloadDefaultCowboysResponse, error)
	GetCowboyByName(ctx context.Context, in *GetCowboyByNameRequest, opts ...grpc.CallOption) (*GetCowboyByNameResponse, error)
	ShootAtRandom(ctx context.Context, in *ShootAtRandomRequest, opts ...grpc.CallOption) (*ShootAtRandomResponse, error)
	Test(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestResponse, error)
}

type cowboysServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCowboysServiceClient(cc grpc.ClientConnInterface) CowboysServiceClient {
	return &cowboysServiceClient{cc}
}

func (c *cowboysServiceClient) Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error) {
	out := new(RunResponse)
	err := c.cc.Invoke(ctx, "/cowboys.CowboysService/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cowboysServiceClient) ReloadDefaultCowboys(ctx context.Context, in *ReloadDefaultCowboysRequest, opts ...grpc.CallOption) (*ReloadDefaultCowboysResponse, error) {
	out := new(ReloadDefaultCowboysResponse)
	err := c.cc.Invoke(ctx, "/cowboys.CowboysService/ReloadDefaultCowboys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cowboysServiceClient) GetCowboyByName(ctx context.Context, in *GetCowboyByNameRequest, opts ...grpc.CallOption) (*GetCowboyByNameResponse, error) {
	out := new(GetCowboyByNameResponse)
	err := c.cc.Invoke(ctx, "/cowboys.CowboysService/GetCowboyByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cowboysServiceClient) ShootAtRandom(ctx context.Context, in *ShootAtRandomRequest, opts ...grpc.CallOption) (*ShootAtRandomResponse, error) {
	out := new(ShootAtRandomResponse)
	err := c.cc.Invoke(ctx, "/cowboys.CowboysService/ShootAtRandom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cowboysServiceClient) Test(ctx context.Context, in *TestRequest, opts ...grpc.CallOption) (*TestResponse, error) {
	out := new(TestResponse)
	err := c.cc.Invoke(ctx, "/cowboys.CowboysService/Test", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CowboysServiceServer is the server API for CowboysService service.
// All implementations must embed UnimplementedCowboysServiceServer
// for forward compatibility
type CowboysServiceServer interface {
	Run(context.Context, *RunRequest) (*RunResponse, error)
	ReloadDefaultCowboys(context.Context, *ReloadDefaultCowboysRequest) (*ReloadDefaultCowboysResponse, error)
	GetCowboyByName(context.Context, *GetCowboyByNameRequest) (*GetCowboyByNameResponse, error)
	ShootAtRandom(context.Context, *ShootAtRandomRequest) (*ShootAtRandomResponse, error)
	Test(context.Context, *TestRequest) (*TestResponse, error)
	mustEmbedUnimplementedCowboysServiceServer()
}

// UnimplementedCowboysServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCowboysServiceServer struct {
}

func (UnimplementedCowboysServiceServer) Run(context.Context, *RunRequest) (*RunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (UnimplementedCowboysServiceServer) ReloadDefaultCowboys(context.Context, *ReloadDefaultCowboysRequest) (*ReloadDefaultCowboysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReloadDefaultCowboys not implemented")
}
func (UnimplementedCowboysServiceServer) GetCowboyByName(context.Context, *GetCowboyByNameRequest) (*GetCowboyByNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCowboyByName not implemented")
}
func (UnimplementedCowboysServiceServer) ShootAtRandom(context.Context, *ShootAtRandomRequest) (*ShootAtRandomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShootAtRandom not implemented")
}
func (UnimplementedCowboysServiceServer) Test(context.Context, *TestRequest) (*TestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Test not implemented")
}
func (UnimplementedCowboysServiceServer) mustEmbedUnimplementedCowboysServiceServer() {}

// UnsafeCowboysServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CowboysServiceServer will
// result in compilation errors.
type UnsafeCowboysServiceServer interface {
	mustEmbedUnimplementedCowboysServiceServer()
}

func RegisterCowboysServiceServer(s grpc.ServiceRegistrar, srv CowboysServiceServer) {
	s.RegisterService(&CowboysService_ServiceDesc, srv)
}

func _CowboysService_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CowboysServiceServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cowboys.CowboysService/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CowboysServiceServer).Run(ctx, req.(*RunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CowboysService_ReloadDefaultCowboys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReloadDefaultCowboysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CowboysServiceServer).ReloadDefaultCowboys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cowboys.CowboysService/ReloadDefaultCowboys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CowboysServiceServer).ReloadDefaultCowboys(ctx, req.(*ReloadDefaultCowboysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CowboysService_GetCowboyByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCowboyByNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CowboysServiceServer).GetCowboyByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cowboys.CowboysService/GetCowboyByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CowboysServiceServer).GetCowboyByName(ctx, req.(*GetCowboyByNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CowboysService_ShootAtRandom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShootAtRandomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CowboysServiceServer).ShootAtRandom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cowboys.CowboysService/ShootAtRandom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CowboysServiceServer).ShootAtRandom(ctx, req.(*ShootAtRandomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CowboysService_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CowboysServiceServer).Test(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cowboys.CowboysService/Test",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CowboysServiceServer).Test(ctx, req.(*TestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CowboysService_ServiceDesc is the grpc.ServiceDesc for CowboysService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CowboysService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cowboys.CowboysService",
	HandlerType: (*CowboysServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Run",
			Handler:    _CowboysService_Run_Handler,
		},
		{
			MethodName: "ReloadDefaultCowboys",
			Handler:    _CowboysService_ReloadDefaultCowboys_Handler,
		},
		{
			MethodName: "GetCowboyByName",
			Handler:    _CowboysService_GetCowboyByName_Handler,
		},
		{
			MethodName: "ShootAtRandom",
			Handler:    _CowboysService_ShootAtRandom_Handler,
		},
		{
			MethodName: "Test",
			Handler:    _CowboysService_Test_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/sumlookup/cowboys/pb/cowboys.proto",
}