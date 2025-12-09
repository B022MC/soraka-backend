package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRootRouter,
	NewLcuRouter,
	NewClientRouter,
	NewMatchRouter,
	NewAssetsRouter,
)

type RootRouter struct {
	lcuRouter    *LcuRouter
	clientRouter *ClientRouter
	matchRouter  *MatchRouter
	assetsRouter *AssetsRouter
}

func (g *RootRouter) InitRouter(group *gin.RouterGroup) {
	g.lcuRouter.InitRouter(group)
	g.clientRouter.InitRouter(group)
	g.matchRouter.InitRouter(group)
	g.assetsRouter.InitRouter(group)
}

func NewRootRouter(
	lcuRouter *LcuRouter,
	clientRouter *ClientRouter,
	matchRouter *MatchRouter,
	assetsRouter *AssetsRouter,
) *RootRouter {
	return &RootRouter{
		lcuRouter:    lcuRouter,
		clientRouter: clientRouter,
		matchRouter:  matchRouter,
		assetsRouter: assetsRouter,
	}
}
