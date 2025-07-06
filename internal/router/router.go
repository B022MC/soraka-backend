package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRootRouter,
	NewLcuRouter,
	NewClientRouter,
)

type RootRouter struct {
	lcuRouter    *LcuRouter
	clientRouter *ClientRouter
}

func (g *RootRouter) InitRouter(group *gin.RouterGroup) {
	g.lcuRouter.InitRouter(group)
	g.clientRouter.InitRouter(group)
}

func NewRootRouter(
	lcuRouter *LcuRouter,
	clientRouter *ClientRouter,
) *RootRouter {
	return &RootRouter{
		lcuRouter:    lcuRouter,
		clientRouter: clientRouter,
	}
}
