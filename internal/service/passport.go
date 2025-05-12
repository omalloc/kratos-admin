package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	adminpb "github.com/omalloc/kratos-admin/api/console/administration"
	pb "github.com/omalloc/kratos-admin/api/console/passport"
	"github.com/omalloc/kratos-admin/internal/biz"
	"github.com/omalloc/kratos-admin/internal/conf"
	"github.com/omalloc/kratos-admin/pkg/jwt"
	"github.com/omalloc/kratos-admin/pkg/tokener"
)

type PassportService struct {
	pb.UnimplementedPassportServer

	userUsecase *biz.UserUsecase
	tokener     tokener.AppToken
}

func NewPassportService(c *conf.Bootstrap, userUsecase *biz.UserUsecase) *PassportService {
	return &PassportService{
		userUsecase: userUsecase,
		tokener: tokener.NewTokener(
			tokener.WithTTL(time.Hour),
			tokener.WithSecret(c.Passport.Secret),
		),
	}
}

func (s *PassportService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := s.userUsecase.Login(ctx, req.Username, req.Password, req.AutoLogin)
	if err != nil {
		return nil, err
	}

	token, err := s.tokener.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	tr, ok := transport.FromServerContext(ctx)
	if ok {
		tr.ReplyHeader().Add("Authorization", token)
	}

	log.Infof("audit -> login user: %+v", user)
	return &pb.LoginReply{}, nil
}
func (s *PassportService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	return &pb.LogoutReply{}, nil
}
func (s *PassportService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *PassportService) SendCaptcha(ctx context.Context, req *pb.SendCaptchaRequest) (*pb.SendCaptchaReply, error) {
	return &pb.SendCaptchaReply{}, nil
}
func (s *PassportService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordReply, error) {
	return &pb.ResetPasswordReply{}, nil
}
func (s *PassportService) CurrentUser(ctx context.Context, req *pb.CurrentUserRequest) (*pb.CurrentUserReply, error) {
	claims, _ := jwt.FromContext(ctx)

	user, err := s.userUsecase.GetUser(ctx, claims.UID)
	if err != nil {
		return nil, err
	}

	return &pb.CurrentUserReply{
		User: &adminpb.UserInfo{
			Id:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Status:    adminpb.UserStatus(user.Status),
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Roles: lo.Map(user.Roles, func(item *biz.Role, _ int) *adminpb.RoleInfo {
			return &adminpb.RoleInfo{
				Id:       item.ID,
				Name:     item.Name,
				Describe: item.Describe,
				Status:   int32(item.Status),
				Permissions: lo.Map(item.Permissions, func(item *biz.RolePermission, _ int) *adminpb.RolePermission {
					return &adminpb.RolePermission{
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
