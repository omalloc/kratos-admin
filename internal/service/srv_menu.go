package service

import (
	"context"

	"github.com/omalloc/contrib/protobuf"
	"github.com/omalloc/kratos-admin/pkg/idgen"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/omalloc/kratos-admin/internal/biz"
)

type MenuService struct {
	pb.UnimplementedMenuServer

	usecase *biz.MenuUsecase
}

func NewMenuService(usecase *biz.MenuUsecase) *MenuService {
	return &MenuService{usecase: usecase}
}

func (s *MenuService) CreateMenu(ctx context.Context, req *pb.CreateMenuRequest) (*pb.CreateMenuReply, error) {
	m := &biz.Menu{
		UID:    idgen.NextId(),
		PID:    req.Pid,
		Name:   req.Name,
		Icon:   req.Icon,
		Path:   req.Path,
		SortBy: req.SortBy,
		Hidden: req.Hidden,
		Status: int32(req.Status),
	}

	if err := s.usecase.Create(ctx, m); err != nil {
		return nil, err
	}

	return &pb.CreateMenuReply{
		Uid: m.UID,
	}, nil
}

func (s *MenuService) UpdateMenu(ctx context.Context, req *pb.UpdateMenuRequest) (*pb.UpdateMenuReply, error) {
	m := &biz.Menu{
		UID:          req.Uid,
		PID:          req.Pid,
		Name:         req.Name,
		Icon:         req.Icon,
		Path:         req.Path,
		SortBy:       req.SortBy,
		Hidden:       req.Hidden,
		PermissionID: req.PermissionId,
		Status:       int32(req.Status),
	}

	if err := s.usecase.Update(ctx, m); err != nil {
		return nil, err
	}

	return &pb.UpdateMenuReply{}, nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, req *pb.DeleteMenuRequest) (*pb.DeleteMenuReply, error) {
	if err := s.usecase.Delete(ctx, req.Uid); err != nil {
		return nil, err
	}

	return &pb.DeleteMenuReply{}, nil
}

func (s *MenuService) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.GetMenuReply, error) {
	m, err := s.usecase.SelectByID(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetMenuReply{
		Data: toMenuProto(m, 0),
	}, nil
}

func (s *MenuService) ListMenu(ctx context.Context, req *pb.ListMenuRequest) (*pb.ListMenuReply, error) {
	pagination := protobuf.PageWrap(req.Pagination)
	menus, err := s.usecase.SelectList(ctx, pagination, req.Name, int32(req.Status))
	if err != nil {
		return nil, err
	}

	return &pb.ListMenuReply{
		Data:       lo.Map(menus, toMenuProto),
		Pagination: pagination.Resp(),
	}, nil
}

func toMenuProto(m *biz.Menu, _ int) *pb.MenuInfo {
	return &pb.MenuInfo{
		Uid:          m.UID,
		Pid:          m.PID,
		PermissionId: m.PermissionID,
		Name:         m.Name,
		Icon:         m.Icon,
		Path:         m.Path,
		SortBy:       m.SortBy,
		Hidden:       m.Hidden,
		Status:       pb.MenuStatus(m.Status),
		CreatedAt:    timestamppb.New(m.CreatedAt),
		UpdatedAt:    timestamppb.New(m.UpdatedAt),
	}
}
