package current_summoner

import (
	currentSummonerRepo "github.com/B022MC/soraka-backend/internal/dal/repo/current_summoner"
	"github.com/go-kratos/kratos/v2/log"
)

type CurrentSummonerUseCase struct {
	repo currentSummonerRepo.CurrentSummonerRepo
	log  *log.Helper
}

func NewCurrentSummonerUseCase(repo currentSummonerRepo.CurrentSummonerRepo, logger log.Logger) *CurrentSummonerUseCase {
	return &CurrentSummonerUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "use/current_summoner"))}
}
