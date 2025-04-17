// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: console/administration/permission.proto

package administration

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Permission_CreatePermission_FullMethodName = "/api.console.administration.Permission/CreatePermission"
	Permission_UpdatePermission_FullMethodName = "/api.console.administration.Permission/UpdatePermission"
	Permission_DeletePermission_FullMethodName = "/api.console.administration.Permission/DeletePermission"
	Permission_GetPermission_FullMethodName    = "/api.console.administration.Permission/GetPermission"
	Permission_ListPermission_FullMethodName   = "/api.console.administration.Permission/ListPermission"
)

// PermissionClient is the client API for Permission service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PermissionClient interface {
	CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*CreatePermissionReply, error)
	UpdatePermission(ctx context.Context, in *UpdatePermissionRequest, opts ...grpc.CallOption) (*UpdatePermissionReply, error)
	DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*DeletePermissionReply, error)
	GetPermission(ctx context.Context, in *GetPermissionRequest, opts ...grpc.CallOption) (*GetPermissionReply, error)
	ListPermission(ctx context.Context, in *ListPermissionRequest, opts ...grpc.CallOption) (*ListPermissionReply, error)
}

type permissionClient struct {
	cc grpc.ClientConnInterface
}

func NewPermissionClient(cc grpc.ClientConnInterface) PermissionClient {
	return &permissionClient{cc}
}

func (c *permissionClient) CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...grpc.CallOption) (*CreatePermissionReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePermissionReply)
	err := c.cc.Invoke(ctx, Permission_CreatePermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) UpdatePermission(ctx context.Context, in *UpdatePermissionRequest, opts ...grpc.CallOption) (*UpdatePermissionReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePermissionReply)
	err := c.cc.Invoke(ctx, Permission_UpdatePermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...grpc.CallOption) (*DeletePermissionReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePermissionReply)
	err := c.cc.Invoke(ctx, Permission_DeletePermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) GetPermission(ctx context.Context, in *GetPermissionRequest, opts ...grpc.CallOption) (*GetPermissionReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPermissionReply)
	err := c.cc.Invoke(ctx, Permission_GetPermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *permissionClient) ListPermission(ctx context.Context, in *ListPermissionRequest, opts ...grpc.CallOption) (*ListPermissionReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPermissionReply)
	err := c.cc.Invoke(ctx, Permission_ListPermission_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PermissionServer is the server API for Permission service.
// All implementations must embed UnimplementedPermissionServer
// for forward compatibility.
type PermissionServer interface {
	CreatePermission(context.Context, *CreatePermissionRequest) (*CreatePermissionReply, error)
	UpdatePermission(context.Context, *UpdatePermissionRequest) (*UpdatePermissionReply, error)
	DeletePermission(context.Context, *DeletePermissionRequest) (*DeletePermissionReply, error)
	GetPermission(context.Context, *GetPermissionRequest) (*GetPermissionReply, error)
	ListPermission(context.Context, *ListPermissionRequest) (*ListPermissionReply, error)
	mustEmbedUnimplementedPermissionServer()
}

// UnimplementedPermissionServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPermissionServer struct{}

func (UnimplementedPermissionServer) CreatePermission(context.Context, *CreatePermissionRequest) (*CreatePermissionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePermission not implemented")
}
func (UnimplementedPermissionServer) UpdatePermission(context.Context, *UpdatePermissionRequest) (*UpdatePermissionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePermission not implemented")
}
func (UnimplementedPermissionServer) DeletePermission(context.Context, *DeletePermissionRequest) (*DeletePermissionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePermission not implemented")
}
func (UnimplementedPermissionServer) GetPermission(context.Context, *GetPermissionRequest) (*GetPermissionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPermission not implemented")
}
func (UnimplementedPermissionServer) ListPermission(context.Context, *ListPermissionRequest) (*ListPermissionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPermission not implemented")
}
func (UnimplementedPermissionServer) mustEmbedUnimplementedPermissionServer() {}
func (UnimplementedPermissionServer) testEmbeddedByValue()                    {}

// UnsafePermissionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PermissionServer will
// result in compilation errors.
type UnsafePermissionServer interface {
	mustEmbedUnimplementedPermissionServer()
}

func RegisterPermissionServer(s grpc.ServiceRegistrar, srv PermissionServer) {
	// If the following call pancis, it indicates UnimplementedPermissionServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Permission_ServiceDesc, srv)
}

func _Permission_CreatePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).CreatePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Permission_CreatePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).CreatePermission(ctx, req.(*CreatePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_UpdatePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).UpdatePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Permission_UpdatePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).UpdatePermission(ctx, req.(*UpdatePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_DeletePermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).DeletePermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Permission_DeletePermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).DeletePermission(ctx, req.(*DeletePermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_GetPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).GetPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Permission_GetPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).GetPermission(ctx, req.(*GetPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Permission_ListPermission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPermissionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PermissionServer).ListPermission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Permission_ListPermission_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PermissionServer).ListPermission(ctx, req.(*ListPermissionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Permission_ServiceDesc is the grpc.ServiceDesc for Permission service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Permission_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.console.administration.Permission",
	HandlerType: (*PermissionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePermission",
			Handler:    _Permission_CreatePermission_Handler,
		},
		{
			MethodName: "UpdatePermission",
			Handler:    _Permission_UpdatePermission_Handler,
		},
		{
			MethodName: "DeletePermission",
			Handler:    _Permission_DeletePermission_Handler,
		},
		{
			MethodName: "GetPermission",
			Handler:    _Permission_GetPermission_Handler,
		},
		{
			MethodName: "ListPermission",
			Handler:    _Permission_ListPermission_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "console/administration/permission.proto",
}
