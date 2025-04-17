package service

import (
	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/omalloc/kratos-admin/api/console"
	"github.com/omalloc/kratos-agent/api/agent"
)

type ConsoleService struct {
	pb.UnimplementedConsoleServer

	log    *log.Helper
	client agent.AgentClient
}

func NewConsoleService(logger log.Logger, client agent.AgentClient) *ConsoleService {
	return &ConsoleService{
		log:    log.NewHelper(logger),
		client: client,
	}
}
