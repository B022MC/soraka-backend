package biz

import (
	clientUseCase "github.com/B022MC/soraka-backend/internal/biz/client"
	currentSummonerUseCaseBiz "github.com/B022MC/soraka-backend/internal/biz/current_summoner"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	currentSummonerUseCaseBiz.NewCurrentSummonerUseCase,
	clientUseCase.NewClientUseCase,
)
