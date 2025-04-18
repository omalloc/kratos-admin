// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: console/administration/role.proto

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

const OperationRoleBindPermission = "/api.console.administration.Role/BindPermission"
const OperationRoleCreateRole = "/api.console.administration.Role/CreateRole"
const OperationRoleDeleteRole = "/api.console.administration.Role/DeleteRole"
const OperationRoleGetRole = "/api.console.administration.Role/GetRole"
const OperationRoleListRole = "/api.console.administration.Role/ListRole"
const OperationRoleUnbindPermission = "/api.console.administration.Role/UnbindPermission"
const OperationRoleUpdateRole = "/api.console.administration.Role/UpdateRole"

type RoleHTTPServer interface {
	BindPermission(context.Context, *BindPermissionRequest) (*BindPermissionReply, error)
	CreateRole(context.Context, *CreateRoleRequest) (*CreateRoleReply, error)
	DeleteRole(context.Context, *DeleteRoleRequest) (*DeleteRoleReply, error)
	GetRole(context.Context, *GetRoleRequest) (*GetRoleReply, error)
	ListRole(context.Context, *ListRoleRequest) (*ListRoleReply, error)
	UnbindPermission(context.Context, *UnbindPermissionRequest) (*UnbindPermissionReply, error)
	UpdateRole(context.Context, *UpdateRoleRequest) (*UpdateRoleReply, error)
}

func RegisterRoleHTTPServer(s *http.Server, srv RoleHTTPServer) {
	r := s.Route("/")
	r.POST("/api/console/role", _Role_CreateRole0_HTTP_Handler(srv))
	r.PUT("/api/console/role/{id}", _Role_UpdateRole0_HTTP_Handler(srv))
	r.DELETE("/api/console/role/{id}", _Role_DeleteRole0_HTTP_Handler(srv))
	r.GET("/api/console/role/{id}", _Role_GetRole0_HTTP_Handler(srv))
	r.GET("/api/console/role", _Role_ListRole0_HTTP_Handler(srv))
	r.PUT("/api/console/role/{id}/permission", _Role_BindPermission0_HTTP_Handler(srv))
	r.PUT("/api/console/role/{id}/permission/{permission_id}", _Role_UnbindPermission0_HTTP_Handler(srv))
}

func _Role_CreateRole0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleCreateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRole(ctx, req.(*CreateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateRoleReply)
		return ctx.Result(200, reply)
	}
}

func _Role_UpdateRole0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRoleRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleUpdateRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRole(ctx, req.(*UpdateRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateRoleReply)
		return ctx.Result(200, reply)
	}
}

func _Role_DeleteRole0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRoleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleDeleteRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteRole(ctx, req.(*DeleteRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteRoleReply)
		return ctx.Result(200, reply)
	}
}

func _Role_GetRole0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRoleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleGetRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRole(ctx, req.(*GetRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRoleReply)
		return ctx.Result(200, reply)
	}
}

func _Role_ListRole0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRoleRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleListRole)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListRole(ctx, req.(*ListRoleRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListRoleReply)
		return ctx.Result(200, reply)
	}
}

func _Role_BindPermission0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in BindPermissionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleBindPermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BindPermission(ctx, req.(*BindPermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BindPermissionReply)
		return ctx.Result(200, reply)
	}
}

func _Role_UnbindPermission0_HTTP_Handler(srv RoleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UnbindPermissionRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoleUnbindPermission)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UnbindPermission(ctx, req.(*UnbindPermissionRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UnbindPermissionReply)
		return ctx.Result(200, reply)
	}
}

type RoleHTTPClient interface {
	BindPermission(ctx context.Context, req *BindPermissionRequest, opts ...http.CallOption) (rsp *BindPermissionReply, err error)
	CreateRole(ctx context.Context, req *CreateRoleRequest, opts ...http.CallOption) (rsp *CreateRoleReply, err error)
	DeleteRole(ctx context.Context, req *DeleteRoleRequest, opts ...http.CallOption) (rsp *DeleteRoleReply, err error)
	GetRole(ctx context.Context, req *GetRoleRequest, opts ...http.CallOption) (rsp *GetRoleReply, err error)
	ListRole(ctx context.Context, req *ListRoleRequest, opts ...http.CallOption) (rsp *ListRoleReply, err error)
	UnbindPermission(ctx context.Context, req *UnbindPermissionRequest, opts ...http.CallOption) (rsp *UnbindPermissionReply, err error)
	UpdateRole(ctx context.Context, req *UpdateRoleRequest, opts ...http.CallOption) (rsp *UpdateRoleReply, err error)
}

type RoleHTTPClientImpl struct {
	cc *http.Client
}

func NewRoleHTTPClient(client *http.Client) RoleHTTPClient {
	return &RoleHTTPClientImpl{client}
}

func (c *RoleHTTPClientImpl) BindPermission(ctx context.Context, in *BindPermissionRequest, opts ...http.CallOption) (*BindPermissionReply, error) {
	var out BindPermissionReply
	pattern := "/api/console/role/{id}/permission"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleBindPermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoleHTTPClientImpl) CreateRole(ctx context.Context, in *CreateRoleRequest, opts ...http.CallOption) (*CreateRoleReply, error) {
	var out CreateRoleReply
	pattern := "/api/console/role"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleCreateRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoleHTTPClientImpl) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...http.CallOption) (*DeleteRoleReply, error) {
	var out DeleteRoleReply
	pattern := "/api/console/role/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleDeleteRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoleHTTPClientImpl) GetRole(ctx context.Context, in *GetRoleRequest, opts ...http.CallOption) (*GetRoleReply, error) {
	var out GetRoleReply
	pattern := "/api/console/role/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleGetRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoleHTTPClientImpl) ListRole(ctx context.Context, in *ListRoleRequest, opts ...http.CallOption) (*ListRoleReply, error) {
	var out ListRoleReply
	pattern := "/api/console/role"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoleListRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoleHTTPClientImpl) UnbindPermission(ctx context.Context, in *UnbindPermissionRequest, opts ...http.CallOption) (*UnbindPermissionReply, error) {
	var out UnbindPermissionReply
	pattern := "/api/console/role/{id}/permission/{permission_id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleUnbindPermission))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoleHTTPClientImpl) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...http.CallOption) (*UpdateRoleReply, error) {
	var out UpdateRoleReply
	pattern := "/api/console/role/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRoleUpdateRole))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
