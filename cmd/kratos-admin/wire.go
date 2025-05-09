//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"github.com/omalloc/kratos-admin/internal/biz"
	"github.com/omalloc/kratos-admin/internal/conf"
	"github.com/omalloc/kratos-admin/internal/data"
	"github.com/omalloc/kratos-admin/internal/discovery"
	"github.com/omalloc/kratos-admin/internal/server"
	"github.com/omalloc/kratos-admin/internal/service"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, *conf.Server, *conf.Data, *conf.Passport, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, discovery.ProviderSet, newApp))
}
