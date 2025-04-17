package discovery

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/omalloc/kratos-agent/api/agent"
)

func NewAgentService(logger log.Logger, dis registry.Discovery) (agent.AgentClient, error) {
	log.Infof("begin connection to grpc discovery kratos-agent")

	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///kratos-agent"),
		grpc.WithDiscovery(dis),
		grpc.WithHealthCheck(false),
		grpc.WithSubset(0),
		grpc.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
			tracing.Client(),
			logging.Client(logger),
		),
	)

	// example ignore err.
	if err != nil {
		return nil, err
	}

	log.Infof("connection to grpc discovery kratos-agent success")
	return agent.NewAgentClient(conn), nil
}
