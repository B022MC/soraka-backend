package biz

import (
	"github.com/google/wire"
	currentSummonerUseCaseBiz "soraka-backend/internal/biz/current_summoner"
)

var ProviderSet = wire.NewSet(
	currentSummonerUseCaseBiz.NewCurrentSummonerUseCase,
)
