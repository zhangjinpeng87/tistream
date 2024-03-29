// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.8.0
// source: dispatcherpb.proto

package __

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

// DispatcherServiceClient is the client API for DispatcherService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DispatcherServiceClient interface {
	NotifiyTenantHasUpdate(ctx context.Context, in *HasNewChangeReq, opts ...grpc.CallOption) (*HasNewChangeResp, error)
	ScheduleNewTenant(ctx context.Context, in *TenantTasksReq, opts ...grpc.CallOption) (*TenantTasksResp, error)
}

type dispatcherServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDispatcherServiceClient(cc grpc.ClientConnInterface) DispatcherServiceClient {
	return &dispatcherServiceClient{cc}
}

func (c *dispatcherServiceClient) NotifiyTenantHasUpdate(ctx context.Context, in *HasNewChangeReq, opts ...grpc.CallOption) (*HasNewChangeResp, error) {
	out := new(HasNewChangeResp)
	err := c.cc.Invoke(ctx, "/tistreampb.DispatcherService/NotifiyTenantHasUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dispatcherServiceClient) ScheduleNewTenant(ctx context.Context, in *TenantTasksReq, opts ...grpc.CallOption) (*TenantTasksResp, error) {
	out := new(TenantTasksResp)
	err := c.cc.Invoke(ctx, "/tistreampb.DispatcherService/ScheduleNewTenant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DispatcherServiceServer is the server API for DispatcherService service.
// All implementations must embed UnimplementedDispatcherServiceServer
// for forward compatibility
type DispatcherServiceServer interface {
	NotifiyTenantHasUpdate(context.Context, *HasNewChangeReq) (*HasNewChangeResp, error)
	ScheduleNewTenant(context.Context, *TenantTasksReq) (*TenantTasksResp, error)
	mustEmbedUnimplementedDispatcherServiceServer()
}

// UnimplementedDispatcherServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDispatcherServiceServer struct {
}

func (UnimplementedDispatcherServiceServer) NotifiyTenantHasUpdate(context.Context, *HasNewChangeReq) (*HasNewChangeResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifiyTenantHasUpdate not implemented")
}
func (UnimplementedDispatcherServiceServer) ScheduleNewTenant(context.Context, *TenantTasksReq) (*TenantTasksResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleNewTenant not implemented")
}
func (UnimplementedDispatcherServiceServer) mustEmbedUnimplementedDispatcherServiceServer() {}

// UnsafeDispatcherServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DispatcherServiceServer will
// result in compilation errors.
type UnsafeDispatcherServiceServer interface {
	mustEmbedUnimplementedDispatcherServiceServer()
}

func RegisterDispatcherServiceServer(s grpc.ServiceRegistrar, srv DispatcherServiceServer) {
	s.RegisterService(&DispatcherService_ServiceDesc, srv)
}

func _DispatcherService_NotifiyTenantHasUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HasNewChangeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatcherServiceServer).NotifiyTenantHasUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tistreampb.DispatcherService/NotifiyTenantHasUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatcherServiceServer).NotifiyTenantHasUpdate(ctx, req.(*HasNewChangeReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DispatcherService_ScheduleNewTenant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TenantTasksReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatcherServiceServer).ScheduleNewTenant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tistreampb.DispatcherService/ScheduleNewTenant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatcherServiceServer).ScheduleNewTenant(ctx, req.(*TenantTasksReq))
	}
	return interceptor(ctx, in, info, handler)
}

// DispatcherService_ServiceDesc is the grpc.ServiceDesc for DispatcherService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DispatcherService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tistreampb.DispatcherService",
	HandlerType: (*DispatcherServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NotifiyTenantHasUpdate",
			Handler:    _DispatcherService_NotifiyTenantHasUpdate_Handler,
		},
		{
			MethodName: "ScheduleNewTenant",
			Handler:    _DispatcherService_ScheduleNewTenant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dispatcherpb.proto",
}
