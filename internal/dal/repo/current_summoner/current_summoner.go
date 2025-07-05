package current_summoner

import (
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type CurrentSummonerRepo interface {
}

type currentSummonerRepo struct {
	data *infra.Data
	log  *log.Helper
}

func NewCurrentSummonerRepo(data *infra.Data, logger log.Logger) CurrentSummonerRepo {
	return &currentSummonerRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/current_summoner")),
	}
}
