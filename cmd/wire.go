//go:build wireinject
// +build wireinject

package main

import (
	"github.com/B022MC/soraka-backend/internal/biz"
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/repo"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/B022MC/soraka-backend/internal/router"
	"github.com/B022MC/soraka-backend/internal/server"
	"github.com/B022MC/soraka-backend/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Global, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {

	panic(wire.Build(server.ProviderSet, infra.ProviderSet, repo.ProviderSet, biz.ProviderSet, service.ProviderSet, router.ProviderSet, newApp))
}
