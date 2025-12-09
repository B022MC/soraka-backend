package repo

import (
	champSelectRepo "github.com/B022MC/soraka-backend/internal/dal/repo/champ_select"
	clientRepo "github.com/B022MC/soraka-backend/internal/dal/repo/client"
	currentSummonerRepo "github.com/B022MC/soraka-backend/internal/dal/repo/current_summoner"
	gameflowRepo "github.com/B022MC/soraka-backend/internal/dal/repo/gameflow"
	lobbyRepo "github.com/B022MC/soraka-backend/internal/dal/repo/lobby"
	matchRepo "github.com/B022MC/soraka-backend/internal/dal/repo/match"
	profileRepo "github.com/B022MC/soraka-backend/internal/dal/repo/profile"
	rankRepo "github.com/B022MC/soraka-backend/internal/dal/repo/rank"
	runesRepo "github.com/B022MC/soraka-backend/internal/dal/repo/runes"
	spectateRepo "github.com/B022MC/soraka-backend/internal/dal/repo/spectate"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	currentSummonerRepo.NewCurrentSummonerRepo,
	clientRepo.NewClientInfoRepo,
	rankRepo.NewRankRepo,
	matchRepo.NewMatchHistoryRepo,
	matchRepo.NewGameDetailRepo,
	// Seraphine 新功能
	gameflowRepo.NewGameflowRepo,
	champSelectRepo.NewChampSelectRepo,
	runesRepo.NewRunesRepo,
	profileRepo.NewProfileRepo,
	lobbyRepo.NewLobbyRepo,
	spectateRepo.NewSpectateRepo,
)
