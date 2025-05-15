package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/samber/lo"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/types/known/timestamppb"

	adminpb "github.com/omalloc/kratos-admin/api/console/administration"
	pb "github.com/omalloc/kratos-admin/api/console/passport"
	"github.com/omalloc/kratos-admin/internal/biz"
	"github.com/omalloc/kratos-admin/internal/conf"
	"github.com/omalloc/kratos-admin/internal/event"
	"github.com/omalloc/kratos-admin/pkg/jwt"
	"github.com/omalloc/kratos-admin/pkg/tokener"
)

const defaultCaptchaLen = 6

type PassportService struct {
	pb.UnimplementedPassportServer

	userUsecase               *biz.UserUsecase
	tokener                   tokener.AppToken
	etcdClient                *clientv3.Client
	applicationEventPublisher *event.ApplicationEventPublisher
}

func NewPassportService(c *conf.Bootstrap, applicationEventPublisher *event.ApplicationEventPublisher, userUsecase *biz.UserUsecase, etcdClient *clientv3.Client) *PassportService {
	return &PassportService{
		applicationEventPublisher: applicationEventPublisher,
		userUsecase:               userUsecase,
		tokener: tokener.NewTokener(
			tokener.WithTTL(time.Hour),
			tokener.WithSecret(c.Passport.Secret),
		),
		etcdClient: etcdClient,
	}
}

// 登录
func (s *PassportService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := s.userUsecase.Login(ctx, req.Username, req.Password, req.AutoLogin)
	if err != nil {
		s.applicationEventPublisher.Publish(ctx, "passport.login.failed", event.NewMessage(event.NewUUID(), []byte(err.Error())))
		return nil, err
	}

	token, err := s.tokener.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	s.applicationEventPublisher.Publish(ctx, "passport.login.token-generated", event.NewMessage(event.NewUUID(), []byte(token)))
	tr, ok := transport.FromServerContext(ctx)
	if ok {
		tr.ReplyHeader().Add("Authorization", token)
	}

	s.applicationEventPublisher.Publish(ctx, "passport.login.success", event.NewMessage(event.NewUUID(), event.Marshal(user)))
	return &pb.LoginReply{}, nil
}

// 登出
func (s *PassportService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	claims, _ := jwt.FromContext(ctx)
	_, _ = s.etcdClient.Delete(ctx, fmt.Sprintf("/app/auth_token/%d", claims.UID))

	s.applicationEventPublisher.Publish(ctx, "passport.logout.success", event.NewMessage(event.NewUUID(), []byte(fmt.Sprintf("%d", claims.UID))))
	return &pb.LogoutReply{}, nil
}

// 发送验证码 (登录的用户)
func (s *PassportService) SendCaptcha(ctx context.Context, req *pb.SendCaptchaRequest) (*pb.SendCaptchaReply, error) {
	claims, _ := jwt.FromContext(ctx)
	code := generateCaptcha(defaultCaptchaLen)

	switch req.Type {
	case pb.CaptchaType_CAPTCHA_TYPE_SMS:
		// TODO: send sms captcha
	case pb.CaptchaType_CAPTCHA_TYPE_EMAIL:
		// TODO: send email captcha
	}

	key := s.fmtPassportCaptchaKey(claims.UID)
	_, err := s.etcdClient.Put(ctx, key, code, clientv3.WithLease(clientv3.LeaseID(0)))
	if err != nil {
		return nil, err
	}
	// 设置过期时间 5 分钟
	lease, err := s.etcdClient.Grant(ctx, 300)
	if err == nil {
		_, _ = s.etcdClient.Put(ctx, key, code, clientv3.WithLease(lease.ID))
	}

	log.Infof("send captcha to user: %d, code: %s", claims.UID, code)

	return &pb.SendCaptchaReply{}, nil
}

// 注册
func (s *PassportService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	// 假设 req 里有 phone 和 captcha 字段
	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	nickname := generateCaptcha(6) // 游客用户
	s.userUsecase.CreateUser(ctx, &biz.User{
		Username: req.Username,
		Email:    req.Email,
		Nickname: fmt.Sprintf("游客用户%s", nickname),
		Password: req.Password,
	})

	return &pb.RegisterReply{}, nil
}

// 发送重置密码验证码到邮箱或短信SMS
func (s *PassportService) SendResetPassword(ctx context.Context, req *pb.SendResetPasswordCaptchaRequest) (*pb.SendResetPasswordCaptchaReply, error) {
	key := s.fmtPassportResetKey(req.Email)

	code := generateCaptcha(defaultCaptchaLen)
	_, err := s.etcdClient.Put(ctx, key, code)
	if err != nil {
		return nil, err
	}

	// TODO: send captcha to email...
	//

	return &pb.SendResetPasswordCaptchaReply{}, nil
}

// 重置密码
func (s *PassportService) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordReply, error) {
	// 假设 req 里有 phone 和 captcha 字段
	if req.Email == "" || req.Password == "" || req.Token == "" {
		return nil, fmt.Errorf("email, new-password and token are required")
	}

	// match email token
	key := s.fmtPassportResetKey(req.Email)
	resp, err := s.etcdClient.Get(ctx, key)
	if err != nil || len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("captcha expired or not found")
	}
	if string(resp.Kvs[0].Value) != req.Token {
		return nil, fmt.Errorf("captcha error")
	}

	// update user password
	err = s.userUsecase.UpdatePassword(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	// 校验通过，删除验证码
	_, _ = s.etcdClient.Delete(ctx, key)

	return &pb.ResetPasswordReply{}, nil
}

// 更新用户名
func (s *PassportService) UpdateUsername(ctx context.Context, req *pb.UpdateUsernameRequest) (*pb.UpdateUsernameReply, error) {
	return &pb.UpdateUsernameReply{}, nil
}

// 更新用户信息
func (s *PassportService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return &pb.UpdateProfileReply{}, nil
}

// 获取当前用户信息
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

func (s *PassportService) fmtPassportCaptchaKey(uid int64) string {
	return fmt.Sprintf("/app/passport/captcha/%d", uid)
}

func (s *PassportService) fmtPassportResetKey(emailOrPhone string) string {
	return fmt.Sprintf("/app/passport/reset/%s", emailOrPhone)
}

// 生成6位随机数字验证码
func generateCaptcha(n int) string {
	min := int64(1)
	for i := 1; i < n; i++ {
		min *= 10
	}
	max := min*10 - 1
	return strconv.FormatInt(rand.Int63n(max-min+1)+min, 10)
}
