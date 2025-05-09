package server

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"go.etcd.io/etcd/server/v3/embed"
)

var _ transport.Server = (*EmbedEtcdServer)(nil)

type EmbedEtcdServer struct {
	server *embed.Etcd
}

func NewEmbedEtcd() (*EmbedEtcdServer, func(), error) {
	cfg := embed.NewConfig()
	cfg.Dir = filepath.Join(os.TempDir(), "embed.etcd")
	cfg.LogLevel = "error"
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Errorf("start embed etcd failed %v", err)
		return nil, nil, err
	}

	cleanup := func() {
		if e != nil {
			e.Close()
		}
	}

	// wait for embed etcd started
	<-e.Server.ReadyNotify()
	log.Info("embed etcd server started")

	return &EmbedEtcdServer{
		server: e,
	}, cleanup, nil
}

// Start implements transport.Server.
func (e *EmbedEtcdServer) Start(context.Context) error {
	// do not waiting etcd server start
	return nil
}

// Stop implements transport.Server.
func (e *EmbedEtcdServer) Stop(context.Context) error {
	return nil
}
