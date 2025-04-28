package service

import (
	"context"

	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"

	pb "github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/omalloc/kratos-admin/internal/biz"
)

type PermissionService struct {
	pb.UnimplementedPermissionServer

	usecase *biz.PermissionUsecase
}

func NewPermissionService(usecase *biz.PermissionUsecase) *PermissionService {
	return &PermissionService{usecase: usecase}
}

func (s *PermissionService) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.CreatePermissionReply, error) {
	if err := s.usecase.CreatePermission(ctx, &biz.Permission{
		Name:     req.Name,
		Alias:    req.Alias,
		Describe: req.Describe,
		Status:   int64(req.Status),
		Actions:  lo.Map(req.Actions, toAction),
	}); err != nil {
		return nil, err
	}

	return &pb.CreatePermissionReply{}, nil
}
func (s *PermissionService) UpdatePermission(ctx context.Context, req *pb.UpdatePermissionRequest) (*pb.UpdatePermissionReply, error) {
	if err := s.usecase.UpdatePermission(ctx, &biz.Permission{
		ID:       req.Id,
		Name:     req.Name,
		Alias:    req.Alias,
		Describe: req.Describe,
		Status:   int64(req.Status),
		Actions:  lo.Map(req.Actions, toAction),
	}); err != nil {
		return nil, err
	}
	return &pb.UpdatePermissionReply{}, nil
}

func (s *PermissionService) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*pb.DeletePermissionReply, error) {
	return &pb.DeletePermissionReply{}, nil
}
func (s *PermissionService) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.GetPermissionReply, error) {
	return &pb.GetPermissionReply{}, nil
}
func (s *PermissionService) ListPermission(ctx context.Context, req *pb.ListPermissionRequest) (*pb.ListPermissionReply, error) {
	pagination := protobuf.PageWrap(req.Pagination)
	permissions, err := s.usecase.ListPermission(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return &pb.ListPermissionReply{
		Data:       lo.Map(permissions, toPermissionInfo),
		Pagination: pagination.Resp(),
	}, nil
}

func toPermissionInfo(permission *biz.Permission, _ int) *pb.PermissionInfo {
	return &pb.PermissionInfo{
		Id:       permission.ID,
		Name:     permission.Name,
		Alias:    permission.Alias,
		Describe: permission.Describe,
		Actions:  lo.Map(permission.Actions, fromAction),
		Status:   pb.PermissionStatus(permission.Status),
	}
}

func toAction(vo *pb.Action, _ int) *biz.Action {
	return &biz.Action{
		Key:      vo.Key,
		Describe: vo.Describe,
		Checked:  vo.Checked,
	}
}

func fromAction(action *biz.Action, _ int) *pb.Action {
	return &pb.Action{
		Key:      action.Key,
		Describe: action.Describe,
		Checked:  action.Checked,
	}
}
