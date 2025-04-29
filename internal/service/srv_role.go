package service

import (
	"context"

	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"

	pb "github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/omalloc/kratos-admin/internal/biz"
)

type RoleService struct {
	pb.UnimplementedRoleServer

	usecase *biz.RoleUsecase
}

func NewRoleService(usecase *biz.RoleUsecase) *RoleService {
	return &RoleService{
		usecase: usecase,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleReply, error) {
	if err := s.usecase.CreateRole(ctx, &biz.Role{
		Name:     req.Name,
		Describe: req.Describe,
		Status:   int64(req.Status),
	}); err != nil {
		return nil, err
	}
	return &pb.CreateRoleReply{}, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleReply, error) {
	if err := s.usecase.UpdateRole(ctx, &biz.Role{
		ID:       req.Id,
		Name:     req.Name,
		Describe: req.Describe,
		Status:   int64(req.Status),
	}); err != nil {
		return nil, err
	}
	return &pb.UpdateRoleReply{}, nil
}
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleReply, error) {
	if err := s.usecase.DeleteRole(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteRoleReply{}, nil
}
func (s *RoleService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleReply, error) {
	role, err := s.usecase.SelectID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetRoleReply{
		Id:       role.ID,
		Name:     role.Name,
		Alias:    role.Alias,
		Describe: role.Describe,
		Status:   int32(role.Status),
		Permissions: lo.Map(role.Permissions, func(item *biz.RolePermission, _ int) *pb.RolePermission {
			return &pb.RolePermission{
				Id:         item.ID,
				RoleId:     item.RoleID,
				PermId:     item.PermID,
				Actions:    lo.Map(item.Actions, fromAction),
				DataAccess: lo.Map(item.DataAccess, fromAction),
			}
		}),
	}, nil
}
func (s *RoleService) ListRole(ctx context.Context, req *pb.ListRoleRequest) (*pb.ListRoleReply, error) {
	pagination := protobuf.PageWrap(req.Pagination)
	roles, err := s.usecase.ListRole(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return &pb.ListRoleReply{
		Data:       lo.Map(roles, s.toMap),
		Pagination: pagination.Resp(),
	}, nil
}

func (s *RoleService) BindPermission(ctx context.Context, req *pb.BindPermissionRequest) (*pb.BindPermissionReply, error) {
	// convert to *biz.Action
	actions := lo.Map(req.Actions, toAction)
	dataAccess := lo.Map(req.DataAccess, toAction)

	if err := s.usecase.BindPermission(ctx, req.Id, req.PermissionId, actions, dataAccess); err != nil {
		return nil, err
	}

	return &pb.BindPermissionReply{}, nil
}
func (s *RoleService) UnbindPermission(ctx context.Context, req *pb.UnbindPermissionRequest) (*pb.UnbindPermissionReply, error) {
	if err := s.usecase.UnbindPermission(ctx, req.Id, req.PermissionId); err != nil {
		return nil, err
	}
	return &pb.UnbindPermissionReply{}, nil
}

func (s *RoleService) toMap(item *biz.Role, _ int) *pb.RoleInfo {
	return &pb.RoleInfo{
		Id:       item.ID,
		Name:     item.Name,
		Describe: item.Describe,
	}
}
