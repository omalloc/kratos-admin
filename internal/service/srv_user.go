package service

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/omalloc/kratos-admin/internal/biz"
)

var (
	ErrPasswordMismatch = errors.New(400, "re-password mismatch", "两次密码不匹配")
)

type UserService struct {
	pb.UnimplementedUserServer

	log     *log.Helper
	usecase *biz.UserUsecase
}

func NewUserService(usecase *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		log:     log.NewHelper(logger),
		usecase: usecase,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	if req.Password != req.RePassword {
		return nil, ErrPasswordMismatch
	}

	user := &biz.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Nickname: req.Nickname,
	}
	if err := s.usecase.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	roleIDs := lo.Map(req.RoleIds, func(item string, _ int) int64 {
		return lo.Must(strconv.ParseInt(item, 10, 64))
	})
	s.usecase.UpdateRole(ctx, user.ID, roleIDs)
	return &pb.CreateUserReply{
		Id: user.ID,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	if req.Password != "" {
		if req.Password != req.RePassword {
			return nil, ErrPasswordMismatch
		}
	}

	user := &biz.User{
		ID:       req.Id,
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   int64(req.Status),
	}
	if req.Password != "" {
		hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hp)
	}

	if err := s.usecase.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	roleIDs := lo.Map(req.RoleIds, func(item string, _ int) int64 {
		return lo.Must(strconv.ParseInt(item, 10, 64))
	})
	s.usecase.UpdateRole(ctx, user.ID, roleIDs)
	return &pb.UpdateUserReply{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	if err := s.usecase.DeleteUser(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteUserReply{}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.usecase.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		User: &pb.UserInfo{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Status:    pb.UserStatus(user.Status),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Roles: lo.Map(user.Roles, func(item *biz.Role, _ int) *pb.RoleInfo {
			return &pb.RoleInfo{
				Id:       item.ID,
				Name:     item.Name,
				Describe: item.Describe,
				Status:   int32(item.Status),
				Permissions: lo.Map(item.Permissions, func(item *biz.RolePermission, _ int) *pb.RolePermission {
					return &pb.RolePermission{
						Id:         item.ID,
						RoleId:     item.RoleID,
						PermId:     item.PermID,
						Actions:    lo.Map(item.Actions, fromAction),
						DataAccess: lo.Map(item.DataAccess, fromAction),
					}
				}),
			}
		}),
	}, nil
}

func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {

	pagination := protobuf.PageWrap(req.Pagination)
	users, err := s.usecase.ListUser(ctx, pagination)
	if err != nil {
		return nil, err
	}
	return &pb.ListUserReply{
		Pagination: pagination.Resp(),
		Data:       lo.Map(users, toMap),
	}, nil
}

func (s *UserService) BindRole(ctx context.Context, req *pb.BindRoleRequest) (*pb.BindRoleReply, error) {
	if err := s.usecase.BindRole(ctx, req.Id, int(req.RoleId)); err != nil {
		return nil, err
	}
	return &pb.BindRoleReply{}, nil
}

func (s *UserService) UnbindRole(ctx context.Context, req *pb.UnbindRoleRequest) (*pb.UnbindRoleReply, error) {
	if err := s.usecase.UnbindRole(ctx, req.Id, int(req.RoleId)); err != nil {
		return nil, err
	}
	return &pb.UnbindRoleReply{}, nil
}

func (s *UserService) UpdateRole(ctx context.Context, userID int64, roleIDs []int64) (*pb.UpdateRoleReply, error) {
	if err := s.usecase.UpdateRole(ctx, userID, roleIDs); err != nil {
		return nil, err
	}
	return &pb.UpdateRoleReply{}, nil
}

func toMap(user *biz.UserInfo, _ int) *pb.UserInfo {
	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Status:    pb.UserStatus(user.Status),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		RoleIds: lo.Map(user.RoleIDs, func(item int64, _ int) int32 {
			return int32(item)
		}),
	}
}
