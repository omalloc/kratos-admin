package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv5 "github.com/golang-jwt/jwt/v5"

	consolepb "github.com/omalloc/kratos-admin/api/console"
	adminpb "github.com/omalloc/kratos-admin/api/console/administration"
	passportpb "github.com/omalloc/kratos-admin/api/console/passport"
	"github.com/omalloc/kratos-admin/internal/conf"
	"github.com/omalloc/kratos-admin/internal/service"
	"github.com/omalloc/kratos-admin/pkg/jwt"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, passportc *conf.Passport, logger log.Logger,
	console *service.ConsoleService,
	// admin
	user *service.UserService,
	role *service.RoleService,
	permission *service.PermissionService,
	passport *service.PassportService,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			tracing.Server(),
			logging.Server(logger),
			// JWT
			selector.Server(
				jwt.Server(func(token *jwtv5.Token) (any, error) {
					return []byte(passportc.Secret), nil
				}),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	consolepb.RegisterConsoleServer(srv, console)

	adminpb.RegisterUserServer(srv, user)
	adminpb.RegisterRoleServer(srv, role)
	adminpb.RegisterPermissionServer(srv, permission)
	passportpb.RegisterPassportServer(srv, passport)
	return srv
}
