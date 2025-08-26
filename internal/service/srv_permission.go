package service

import (
	"context"

	"github.com/omalloc/contrib/protobuf"
	"github.com/omalloc/kratos-admin/pkg/idgen"
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
		UID:      idgen.NextId(),
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
		UID:      req.Uid,
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
	if err := s.usecase.DeletePermission(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &pb.DeletePermissionReply{}, nil
}

func (s *PermissionService) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.GetPermissionReply, error) {
	permission, err := s.usecase.GetPermission(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetPermissionReply{
		Uid:      permission.UID,
		Name:     permission.Name,
		Alias:    permission.Alias,
		Describe: permission.Describe,
		Actions:  lo.Map(permission.Actions, fromAction),
		Status:   pb.PermissionStatus(permission.Status),
	}, nil
}
func (s *PermissionService) ListPermission(ctx context.Context, req *pb.ListPermissionRequest) (*pb.ListPermissionReply, error) {
	pagination := protobuf.PageWrap(req.Pagination)
	permissions, err := s.usecase.ListPermission(ctx, req.Name, int32(req.Status), pagination)
	if err != nil {
		return nil, err
	}

	return &pb.ListPermissionReply{
		Data:       lo.Map(permissions, s.toMap),
		Pagination: pagination.Resp(),
	}, nil
}

func (s *PermissionService) ListAllPermission(ctx context.Context, req *pb.ListAllPermissionRequest) (*pb.ListAllPermissionReply, error) {
	permissions, err := s.usecase.ListAllPermission(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListAllPermissionReply{
		Data: lo.Map(permissions, s.toMap),
	}, nil
}

func (s *PermissionService) toMap(permission *biz.Permission, _ int) *pb.PermissionInfo {
	return &pb.PermissionInfo{
		Uid:         permission.UID,
		Name:        permission.Name,
		Alias:       permission.Alias,
		Describe:    permission.Describe,
		AllowDelete: len(permission.Tags) <= 0, // 没有 tag 才可以删除
		Actions:     lo.Map(permission.Actions, fromAction),
		Status:      pb.PermissionStatus(permission.Status),
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
