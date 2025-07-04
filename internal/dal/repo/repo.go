package repo

import (
	"github.com/google/wire"
	currentSummonerRepo "soraka-backend/internal/dal/repo/current_summoner"
)

var ProviderSet = wire.NewSet(
	currentSummonerRepo.NewCurrentSummonerRepo,
)
