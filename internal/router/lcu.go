package router

import (
	automationService "github.com/B022MC/soraka-backend/internal/service/automation"
	champSelectService "github.com/B022MC/soraka-backend/internal/service/champ_select"
	gameflowService "github.com/B022MC/soraka-backend/internal/service/gameflow"
	lcuService "github.com/B022MC/soraka-backend/internal/service/lcu"
	"github.com/gin-gonic/gin"
)

type LcuRouter struct {
	currentSummonerService *lcuService.CurrentSummonerService
	auxiliaryService       *lcuService.AuxiliaryService
	gameflowService        *gameflowService.GameflowService
	champSelectService     *champSelectService.ChampSelectService
	automationService      *automationService.AutomationService
}

func (r *LcuRouter) InitRouter(root *gin.RouterGroup) {
	r.currentSummonerService.RegisterRouter(root)
	r.auxiliaryService.RegisterRouter(root)
	r.gameflowService.RegisterRouter(root)
	r.champSelectService.RegisterRouter(root)
	r.automationService.RegisterRouter(root)
}

func NewLcuRouter(
	currentSummonerService *lcuService.CurrentSummonerService,
	auxiliaryService *lcuService.AuxiliaryService,
	gameflowService *gameflowService.GameflowService,
	champSelectService *champSelectService.ChampSelectService,
	automationService *automationService.AutomationService,
) *LcuRouter {
	return &LcuRouter{
		currentSummonerService: currentSummonerService,
		auxiliaryService:       auxiliaryService,
		gameflowService:        gameflowService,
		champSelectService:     champSelectService,
		automationService:      automationService,
	}
}
