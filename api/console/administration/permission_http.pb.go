// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: console/administration/permission.proto

package administration

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationPermissionCreatePermission = "/api.console.administration.Permission/CreatePermission"
const OperationPermissionDeletePermission = "/api.console.administration.Permission/DeletePermission"
const OperationPermissionGetPermission = "/api.console.administration.Permission/GetPermission"
const OperationPermissionListAllPermission = "/api.console.administration.Permission/ListAllPermission"
const OperationPermissionListPermission = "/api.console.administration.Permission/ListPermission"
const OperationPermissionUpdatePermission = "/api.console.administration.Permission/UpdatePermission"

type PermissionHTTPServer interface {
	CreatePermission(context.Context, *CreatePermissionRequest) (*CreatePermissionReply, error)
	DeletePermission(context.Context, *DeletePermissionRequest) (*DeletePermissionReply, error)
	GetPermission(context.Context, *GetPermissionRequest) (*GetPermissionReply, error)
	ListAllPermission(context.Context, *ListAllPermissionRequest) (*ListAllPermissionReply, error)
	ListPermission(context.Context, *ListPermissionRequest) (*ListPermissionReply, error)
	UpdatePermission(context.Context, *UpdatePermissionRequest) (*UpdatePermissionReply, error)
}

func RegisterPermissionHTTPServer(s *http.Server, srv PermissionHTTPServer) {
	r := s.Route("/")
	r.POST("/api/console/permission", _Permission_CreatePermission0_HTTP_Handler(srv))
	r.PUT("/api/console/permission/{id}", _Permission_UpdatePermission0_HTTP_Handler(srv))
	r.DELETE("/api/console/permission/{id}", _Permission_DeletePermission0_HTTP_Handler(srv))
	r.GET("/api/console/permission/{id}", _Permission_GetPermission0_HTTP_Handler(srv))
	r.GET("/api/console/permission", _Permission_ListPermission0_HTTP_Handler(srv))
	r.GET("/api/console/permission-scoped", _Permission_ListAllPermission0_HTTP_Handler(srv))
}

func _Permission_CreatePermission0_HTTP_Handler(srv PermissionHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreatePermissionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPermissionCreatePermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreatePermission(ctx, req.(*CreatePermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreatePermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Permission_UpdatePermission0_HTTP_Handler(srv PermissionHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdatePermissionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPermissionUpdatePermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdatePermission(ctx, req.(*UpdatePermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdatePermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Permission_DeletePermission0_HTTP_Handler(srv PermissionHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeletePermissionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPermissionDeletePermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeletePermission(ctx, req.(*DeletePermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeletePermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Permission_GetPermission0_HTTP_Handler(srv PermissionHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetPermissionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPermissionGetPermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetPermission(ctx, req.(*GetPermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetPermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Permission_ListPermission0_HTTP_Handler(srv PermissionHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListPermissionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPermissionListPermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListPermission(ctx, req.(*ListPermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListPermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Permission_ListAllPermission0_HTTP_Handler(srv PermissionHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListAllPermissionRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPermissionListAllPermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListAllPermission(ctx, req.(*ListAllPermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListAllPermissionReply)
		return ctx.Result(200, reply)
	}
}

type PermissionHTTPClient interface {
	CreatePermission(ctx context.Context, req *CreatePermissionRequest, opts ...http.CallOption) (rsp *CreatePermissionReply, err error)
	DeletePermission(ctx context.Context, req *DeletePermissionRequest, opts ...http.CallOption) (rsp *DeletePermissionReply, err error)
	GetPermission(ctx context.Context, req *GetPermissionRequest, opts ...http.CallOption) (rsp *GetPermissionReply, err error)
	ListAllPermission(ctx context.Context, req *ListAllPermissionRequest, opts ...http.CallOption) (rsp *ListAllPermissionReply, err error)
	ListPermission(ctx context.Context, req *ListPermissionRequest, opts ...http.CallOption) (rsp *ListPermissionReply, err error)
	UpdatePermission(ctx context.Context, req *UpdatePermissionRequest, opts ...http.CallOption) (rsp *UpdatePermissionReply, err error)
}

type PermissionHTTPClientImpl struct {
	cc *http.Client
}

func NewPermissionHTTPClient(client *http.Client) PermissionHTTPClient {
	return &PermissionHTTPClientImpl{client}
}

func (c *PermissionHTTPClientImpl) CreatePermission(ctx context.Context, in *CreatePermissionRequest, opts ...http.CallOption) (*CreatePermissionReply, error) {
	var out CreatePermissionReply
	pattern := "/api/console/permission"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPermissionCreatePermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PermissionHTTPClientImpl) DeletePermission(ctx context.Context, in *DeletePermissionRequest, opts ...http.CallOption) (*DeletePermissionReply, error) {
	var out DeletePermissionReply
	pattern := "/api/console/permission/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPermissionDeletePermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PermissionHTTPClientImpl) GetPermission(ctx context.Context, in *GetPermissionRequest, opts ...http.CallOption) (*GetPermissionReply, error) {
	var out GetPermissionReply
	pattern := "/api/console/permission/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPermissionGetPermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PermissionHTTPClientImpl) ListAllPermission(ctx context.Context, in *ListAllPermissionRequest, opts ...http.CallOption) (*ListAllPermissionReply, error) {
	var out ListAllPermissionReply
	pattern := "/api/console/permission-scoped"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPermissionListAllPermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PermissionHTTPClientImpl) ListPermission(ctx context.Context, in *ListPermissionRequest, opts ...http.CallOption) (*ListPermissionReply, error) {
	var out ListPermissionReply
	pattern := "/api/console/permission"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPermissionListPermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PermissionHTTPClientImpl) UpdatePermission(ctx context.Context, in *UpdatePermissionRequest, opts ...http.CallOption) (*UpdatePermissionReply, error) {
	var out UpdatePermissionReply
	pattern := "/api/console/permission/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPermissionUpdatePermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
