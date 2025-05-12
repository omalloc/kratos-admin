package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	adminpb "github.com/omalloc/kratos-admin/api/console/administration"
	passportpb "github.com/omalloc/kratos-admin/api/console/passport"
	"github.com/omalloc/kratos-admin/internal/conf"
	"github.com/omalloc/kratos-admin/internal/service"
	"github.com/omalloc/kratos-admin/pkg/jwt"
)

func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList[passportpb.OperationPassportLogin] = struct{}{}
	whiteList[passportpb.OperationPassportLogout] = struct{}{}
	whiteList[passportpb.OperationPassportRegister] = struct{}{}
	whiteList[passportpb.OperationPassportResetPassword] = struct{}{}
	whiteList[passportpb.OperationPassportSendCaptcha] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, passportc *conf.Passport, logger log.Logger,
	// admin
	user *service.UserService,
	role *service.RoleService,
	permission *service.PermissionService,
	passport *service.PassportService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			tracing.Server(),
			logging.Server(logger),
			// JWT
			selector.Server(
				jwt.Server(func(token *jwtv5.Token) (any, error) {
					return []byte(passportc.Secret), nil
				}),
			).Match(NewWhiteListMatcher()).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	srv.HandlePrefix("/q/", openapiv2.NewHandler())

	adminpb.RegisterUserHTTPServer(srv, user)
	adminpb.RegisterRoleHTTPServer(srv, role)
	adminpb.RegisterPermissionHTTPServer(srv, permission)
	passportpb.RegisterPassportHTTPServer(srv, passport)
	return srv
}
