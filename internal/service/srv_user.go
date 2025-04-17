package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"

	"github.com/omalloc/contrib/protobuf"
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
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			RoleIds:   lo.Map(user.Roles, func(item *biz.Role, _ int) int64 { return item.ID }),
		},
	}, nil
}

func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	pagination := protobuf.PageWrap(req.Pagination)
	users, err := s.usecase.ListUser(ctx, pagination)
	if err != nil {
		return nil, err
	}

	userInfos := make([]*pb.UserInfo, 0, len(users))
	for _, user := range users {
		userInfos = append(userInfos, &pb.UserInfo{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Status:    pb.UserStatus(user.Status),
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &pb.ListUserReply{
		Pagination: pagination.Resp(),
		Data:       userInfos,
	}, nil
}

func (s *UserService) BindNamespace(ctx context.Context, req *pb.BindNamespaceRequest) (*pb.BindNamespaceReply, error) {
	if err := s.usecase.BindNamespace(ctx, req.Id, req.NamespaceId); err != nil {
		return nil, err
	}
	return &pb.BindNamespaceReply{}, nil
}

func (s *UserService) UnbindNamespace(ctx context.Context, req *pb.UnbindNamespaceRequest) (*pb.UnbindNamespaceReply, error) {
	if err := s.usecase.UnbindNamespace(ctx, req.Id, req.NamespaceId); err != nil {
		return nil, err
	}
	return &pb.UnbindNamespaceReply{}, nil
}

func (s *UserService) BindRole(ctx context.Context, req *pb.BindRoleRequest) (*pb.BindRoleReply, error) {
	if err := s.usecase.BindRole(ctx, req.Id, req.RoleId); err != nil {
		return nil, err
	}
	return &pb.BindRoleReply{}, nil
}

func (s *UserService) UnbindRole(ctx context.Context, req *pb.UnbindRoleRequest) (*pb.UnbindRoleReply, error) {
	if err := s.usecase.UnbindRole(ctx, req.Id, req.RoleId); err != nil {
		return nil, err
	}
	return &pb.UnbindRoleReply{}, nil
}
