package current_summoner

import (
	"github.com/go-kratos/kratos/v2/log"
	"soraka-backend/internal/infra"
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
