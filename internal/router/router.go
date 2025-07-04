package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRootRouter,
	NewLcuRouter,
)

type RootRouter struct {
	lcuRouter *LcuRouter
}

func (g *RootRouter) InitRouter(group *gin.RouterGroup) {

	g.lcuRouter.InitRouter(group)

}

func NewRootRouter(
	lcuRouter *LcuRouter,
) *RootRouter {
	return &RootRouter{
		lcuRouter: lcuRouter,
	}
}
