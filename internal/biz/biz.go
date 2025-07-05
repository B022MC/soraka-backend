package biz

import (
	currentSummonerUseCaseBiz "github.com/B022MC/soraka-backend/internal/biz/current_summoner"
	gamePhaseBiz "github.com/B022MC/soraka-backend/internal/biz/game_phase"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	currentSummonerUseCaseBiz.NewCurrentSummonerUseCase,
	gamePhaseBiz.NewGamePhaseUseCase,
)
