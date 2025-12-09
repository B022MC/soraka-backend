package biz

import (
	automationUseCase "github.com/B022MC/soraka-backend/internal/biz/automation"
	champSelectUseCase "github.com/B022MC/soraka-backend/internal/biz/champ_select"
	clientUseCase "github.com/B022MC/soraka-backend/internal/biz/client"
	currentSummonerUseCaseBiz "github.com/B022MC/soraka-backend/internal/biz/current_summoner"
	gameflowUseCase "github.com/B022MC/soraka-backend/internal/biz/gameflow"
	matchUseCase "github.com/B022MC/soraka-backend/internal/biz/match"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	currentSummonerUseCaseBiz.NewCurrentSummonerUseCase,
	clientUseCase.NewClientUseCase,
	matchUseCase.NewMatchUseCase,
	// Seraphine 新功能
	gameflowUseCase.NewGameflowUseCase,
	champSelectUseCase.NewChampSelectUseCase,
	automationUseCase.NewAutomationUseCase,
)
