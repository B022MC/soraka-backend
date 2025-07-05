package repo

import (
	currentSummonerRepo "github.com/B022MC/soraka-backend/internal/dal/repo/current_summoner"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	currentSummonerRepo.NewCurrentSummonerRepo,
)
