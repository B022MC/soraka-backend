package repo

import (
	clientRepo "github.com/B022MC/soraka-backend/internal/dal/repo/client"
	currentSummonerRepo "github.com/B022MC/soraka-backend/internal/dal/repo/current_summoner"
	matchRepo "github.com/B022MC/soraka-backend/internal/dal/repo/match"
	rankRepo "github.com/B022MC/soraka-backend/internal/dal/repo/rank"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	currentSummonerRepo.NewCurrentSummonerRepo,
	clientRepo.NewClientInfoRepo,
	rankRepo.NewRankRepo,
	matchRepo.NewMatchHistoryRepo,
	matchRepo.NewGameDetailRepo,
)
