package router

import (
	lcuService "github.com/B022MC/soraka-backend/internal/service/lcu"
	"github.com/gin-gonic/gin"
)

type LcuRouter struct {
	currentSummonerService *lcuService.CurrentSummonerService
	gamePhaseService       *lcuService.GamePhaseService
}

func (r *LcuRouter) InitRouter(root *gin.RouterGroup) {
	r.currentSummonerService.RegisterRouter(root)
	r.gamePhaseService.RegisterRouter(root)

}

func NewLcuRouter(
	currentSummonerService *lcuService.CurrentSummonerService,
	gamePhaseService *lcuService.GamePhaseService,
) *LcuRouter {
	return &LcuRouter{
		currentSummonerService: currentSummonerService,
		gamePhaseService:       gamePhaseService,
	}
}
