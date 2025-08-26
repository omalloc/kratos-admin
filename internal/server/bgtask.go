package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/robfig/cron/v3"
)

var _ transport.Server = (*BackgroundTaskManager)(nil)

type BackgroundTaskManager struct {
	cron *cron.Cron
}

func NewBackgroundTaskManager() *BackgroundTaskManager {
	return &BackgroundTaskManager{
		cron: cron.New(cron.WithSeconds()),
	}
}

// Start implements transport.Server.
func (r *BackgroundTaskManager) Start(context.Context) error {
	r.cron.Start()
	return nil
}

// Stop implements transport.Server.
func (r *BackgroundTaskManager) Stop(context.Context) error {
	r.cron.Stop()
	return nil
}
