package game_phase

import (
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/B022MC/soraka-backend/internal/infra/lcu"
	"github.com/go-kratos/kratos/v2/log"
)

type GamePhaseRepo interface {
	Subscribe() chan string
	Unsubscribe(ch chan string)
}

type gamePhaseRepo struct {
	client *lcu.Client
	log    *log.Helper
}

func NewGamePhaseRepo(data *infra.Data, logger log.Logger) GamePhaseRepo {
	return &gamePhaseRepo{
		client: data.LCU,
		log:    log.NewHelper(log.With(logger, "module", "repo/game_phase")),
	}
}

func (r *gamePhaseRepo) Subscribe() chan string {
	return r.client.SubscribePhase()
}

func (r *gamePhaseRepo) Unsubscribe(ch chan string) {
	r.client.UnsubscribePhase(ch)
}
