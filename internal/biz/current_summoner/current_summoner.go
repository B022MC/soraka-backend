package current_summoner

import (
	"github.com/go-kratos/kratos/v2/log"
	currentSummonerRepo "soraka-backend/internal/dal/repo/current_summoner"
)

type CurrentSummonerUseCase struct {
	repo currentSummonerRepo.CurrentSummonerRepo
	log  *log.Helper
}

func NewCurrentSummonerUseCase(repo currentSummonerRepo.CurrentSummonerRepo, logger log.Logger) *CurrentSummonerUseCase {
	return &CurrentSummonerUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "use/current_summoner"))}
}
