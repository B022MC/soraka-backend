package service

import (
	automationService "github.com/B022MC/soraka-backend/internal/service/automation"
	champSelectService "github.com/B022MC/soraka-backend/internal/service/champ_select"
	clientService "github.com/B022MC/soraka-backend/internal/service/client"
	gameflowService "github.com/B022MC/soraka-backend/internal/service/gameflow"
	lcuService "github.com/B022MC/soraka-backend/internal/service/lcu"
	matchService "github.com/B022MC/soraka-backend/internal/service/match"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	lcuService.NewCurrentSummonerService,
	lcuService.NewAuxiliaryService,
	clientService.NewClientInfoService,
	matchService.NewMatchService,
	// Seraphine 新功能
	gameflowService.NewGameflowService,
	champSelectService.NewChampSelectService,
	automationService.NewAutomationService,
)
