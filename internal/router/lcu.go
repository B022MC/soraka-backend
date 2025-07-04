package router

import (
	"github.com/gin-gonic/gin"
	lcuService "soraka-backend/internal/service/lcu"
)

type LcuRouter struct {
	currentSummonerService *lcuService.CurrentSummonerService
}

func (r *LcuRouter) InitRouter(root *gin.RouterGroup) {
	r.currentSummonerService.RegisterRouter(root)

}

func NewLcuRouter(
	currentSummonerService *lcuService.CurrentSummonerService,
) *LcuRouter {
	return &LcuRouter{
		currentSummonerService: currentSummonerService,
	}
}
