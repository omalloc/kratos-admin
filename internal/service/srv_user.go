package service

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/protobuf"
	"github.com/omalloc/kratos-admin/pkg/idgen"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/omalloc/kratos-admin/api/console/administration"
	"github.com/omalloc/kratos-admin/internal/biz"
)

var ErrPasswordMismatch = errors.New(400, "re-password mismatch", "两次密码不匹配")

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
		UID:       idgen.NextId(),
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Nickname:  req.Nickname,
		Status:    int64(req.Status),
		LastLogin: time.Now(),
	}
	if err := s.usecase.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	roleIDs := lo.Map(req.RoleIds, func(item string, _ int) int64 {
		return lo.Must(strconv.ParseInt(item, 10, 64))
	})
	if err := s.usecase.UpdateRole(ctx, user.UID, roleIDs); err != nil {
		return nil, err
	}
	return &pb.CreateUserReply{
		Uid: user.UID,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	user := &biz.User{
		UID:    req.Uid,
		Status: int64(req.Status),
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.Password != "" {
		// re-password check
		if req.Password != req.RePassword {
			return nil, ErrPasswordMismatch
		}

		hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hp)
	}

	if err := s.usecase.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	if len(req.RoleIds) > 0 {
		roleIDs := lo.Map(req.RoleIds, func(item string, _ int) int64 {
			return lo.Must(strconv.ParseInt(item, 10, 64))
		})
		if err := s.usecase.UpdateRole(ctx, user.UID, roleIDs); err != nil {
			return nil, err
		}
	}
	return &pb.UpdateUserReply{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	if err := s.usecase.DeleteUser(ctx, req.Uid); err != nil {
		return nil, err
	}
	return &pb.DeleteUserReply{}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.usecase.GetUser(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		User: &pb.UserInfo{
			Uid:       user.UID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Status:    pb.UserStatus(user.Status),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Roles: lo.Map(user.Roles, func(item *biz.Role, _ int) *pb.RoleInfo {
			return &pb.RoleInfo{
				Uid:      item.UID,
				Name:     item.Name,
				Describe: item.Describe,
				Status:   int32(item.Status),
				Permissions: lo.Map(item.Permissions, func(item *biz.RolePermission, _ int) *pb.RolePermission {
					return &pb.RolePermission{
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
	users, err := s.usecase.ListUser(ctx, pagination, &biz.UserQueryFilter{
		Status:   int(req.Status),
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ListUserReply{
		Pagination: pagination.Resp(),
		Data:       lo.Map(users, toMap),
	}, nil
}

func (s *UserService) BindRole(ctx context.Context, req *pb.BindRoleRequest) (*pb.BindRoleReply, error) {
	if err := s.usecase.BindRole(ctx, req.Uid, req.RoleId); err != nil {
		return nil, err
	}
	return &pb.BindRoleReply{}, nil
}

func (s *UserService) UnbindRole(ctx context.Context, req *pb.UnbindRoleRequest) (*pb.UnbindRoleReply, error) {
	if err := s.usecase.UnbindRole(ctx, req.Uid, req.RoleId); err != nil {
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
		Uid:       user.UID,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Status:    pb.UserStatus(user.Status),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		LastLogin: lo.Ternary(user.LastLogin.IsZero(), nil, timestamppb.New(user.LastLogin)),
		RoleIds:   user.RoleIDs,
	}
}
