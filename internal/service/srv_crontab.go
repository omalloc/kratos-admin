package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/omalloc/kratos-admin/internal/biz"
)

type CrontabService struct {
	pb.UnimplementedCrontabServer

	log     *log.Helper
	usecase *biz.CrontabUsecase
}

func NewCrontabService(usecase *biz.CrontabUsecase, logger log.Logger) *CrontabService {
	return &CrontabService{
		log:     log.NewHelper(logger),
		usecase: usecase,
	}
}

func (s *CrontabService) CreateCrontab(ctx context.Context, req *pb.CreateCrontabRequest) (*pb.CreateCrontabReply, error) {
	crontab := &biz.Crontab{
		Name:     req.Name,
		Expr:     req.Expr,
		Action:   req.Action,
		Describe: req.Describe,
	}

	if err := s.usecase.CreateCrontab(ctx, crontab); err != nil {
		return nil, err
	}

	return &pb.CreateCrontabReply{}, nil
}

func (s *CrontabService) UpdateCrontab(ctx context.Context, req *pb.UpdateCrontabRequest) (*pb.UpdateCrontabReply, error) {
	crontab := &biz.Crontab{
		UID:      req.Uid,
		Name:     req.Name,
		Expr:     req.Expr,
		Action:   req.Action,
		Describe: req.Describe,
	}

	if err := s.usecase.UpdateCrontab(ctx, crontab); err != nil {
		return nil, err
	}

	return &pb.UpdateCrontabReply{}, nil
}

func (s *CrontabService) DeleteCrontab(ctx context.Context, req *pb.DeleteCrontabRequest) (*pb.DeleteCrontabReply, error) {
	if err := s.usecase.DeleteCrontab(ctx, req.Uid); err != nil {
		return nil, err
	}

	return &pb.DeleteCrontabReply{}, nil
}

func (s *CrontabService) GetCrontab(ctx context.Context, req *pb.GetCrontabRequest) (*pb.GetCrontabReply, error) {
	crontab, err := s.usecase.GetCrontab(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetCrontabReply{
		Data: s.toMap(crontab, 0),
	}, nil
}

func (s *CrontabService) ListCrontab(ctx context.Context, req *pb.ListCrontabRequest) (*pb.ListCrontabReply, error) {
	pagination := protobuf.PageWrap(req.Pagination)

	listCrontab, err := s.usecase.ListCrontab(ctx, pagination)
	if err != nil {
		return nil, err
	}

	return &pb.ListCrontabReply{
		Pagination: pagination.Resp(),
		Data:       lo.Map(listCrontab, s.toMap),
	}, nil
}

func (s *CrontabService) toMap(crontab *biz.Crontab, _ int) *pb.CrontabInfo {
	info := &pb.CrontabInfo{
		Uid:       crontab.UID,
		Name:      crontab.Name,
		Expr:      crontab.Expr,
		Action:    crontab.Action,
		Describe:  crontab.Describe,
		CreatedAt: timestamppb.New(crontab.CreatedAt),
		UpdatedAt: timestamppb.New(crontab.UpdatedAt),
	}

	// 处理可能为空的 LastRunAt
	if crontab.LastRunAt != nil {
		info.LastRunAt = timestamppb.New(*crontab.LastRunAt)
	}

	return info
}
