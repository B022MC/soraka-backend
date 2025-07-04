//go:build wireinject
// +build wireinject

package main

import (
	"soraka-backend/internal/biz"
	"soraka-backend/internal/conf"
	"soraka-backend/internal/dal/repo"
	"soraka-backend/internal/infra"
	"soraka-backend/internal/router"
	"soraka-backend/internal/server"
	"soraka-backend/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Global, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {

	panic(wire.Build(server.ProviderSet, infra.ProviderSet, repo.ProviderSet, biz.ProviderSet, service.ProviderSet, router.ProviderSet, newApp))
}
