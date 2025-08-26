package service

import (
	"context"

	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		Alias:    req.Alias,
		Status:   int64(req.Status),
	}); err != nil {
		return nil, err
	}
	return &pb.CreateRoleReply{}, nil
}
func (s *RoleService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleReply, error) {
	if err := s.usecase.UpdateRole(ctx, &biz.Role{
		UID:      req.Uid,
		Name:     req.Name,
		Alias:    req.Alias,
		Describe: req.Describe,
		Status:   int64(req.Status),
	}); err != nil {
		return nil, err
	}
	return &pb.UpdateRoleReply{}, nil
}
func (s *RoleService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleReply, error) {
	if err := s.usecase.DeleteRole(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &pb.DeleteRoleReply{}, nil
}
func (s *RoleService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleReply, error) {
	role, err := s.usecase.SelectID(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return &pb.GetRoleReply{
		Uid:      role.UID,
		Name:     role.Name,
		Alias:    role.Alias,
		Describe: role.Describe,
		Status:   int32(role.Status),
		Permissions: lo.Map(role.Permissions, func(item *biz.RolePermission, _ int) *pb.RolePermission {
			return &pb.RolePermission{
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
		Data:       lo.Map(roles, s.toRoleJoinedMap),
		Pagination: pagination.Resp(),
	}, nil
}

func (s *RoleService) BindPermission(ctx context.Context, req *pb.BindPermissionRequest) (*pb.BindPermissionReply, error) {
	// convert to *biz.Action
	for _, item := range req.Data {
		actions := lo.Map(item.Actions, toAction)
		dataAccess := lo.Map(item.DataAccess, toAction)
		if err := s.usecase.BindPermission(ctx, req.Uid, item.PermissionId, actions, dataAccess); err != nil {
			return nil, err
		}
	}

	return &pb.BindPermissionReply{}, nil
}
func (s *RoleService) UnbindPermission(ctx context.Context, req *pb.UnbindPermissionRequest) (*pb.UnbindPermissionReply, error) {
	if err := s.usecase.UnbindPermission(ctx, req.Uid, req.PermissionId); err != nil {
		return nil, err
	}
	return &pb.UnbindPermissionReply{}, nil
}

func (s *RoleService) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllReply, error) {
	roles, err := s.usecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllReply{
		Data: lo.Map(roles, s.toRoleMap),
	}, nil
}

func (s *RoleService) toRoleMap(item *biz.Role, _ int) *pb.RoleInfo {
	return &pb.RoleInfo{
		Uid:         item.UID,
		Name:        item.Name,
		Describe:    item.Describe,
		Alias:       item.Alias,
		Status:      int32(item.Status),
		Permissions: lo.Map(item.Permissions, s.toPermission),
	}
}

func (s *RoleService) toRoleJoinedMap(item *biz.RoleJoinPermission, _ int) *pb.RoleInfo {
	return &pb.RoleInfo{
		Uid:         item.UID,
		Name:        item.Name,
		Describe:    item.Describe,
		Alias:       item.Alias,
		Status:      int32(item.Status),
		Permissions: lo.Map(item.Permissions, s.toPermission),
	}
}
func (s *RoleService) toPermission(item *biz.RolePermission, _ int) *pb.RolePermission {
	return &pb.RolePermission{
		RoleId:     item.RoleID,
		PermId:     item.PermID,
		Actions:    lo.Map(item.Actions, fromAction),
		DataAccess: lo.Map(item.DataAccess, fromAction),
		CreatedAt:  timestamppb.New(item.CreatedAt),
	}
}
