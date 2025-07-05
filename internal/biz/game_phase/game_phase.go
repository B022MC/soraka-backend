package game_phase

import (
	repo "github.com/B022MC/soraka-backend/internal/dal/repo/game_phase"
	"github.com/go-kratos/kratos/v2/log"
)

type GamePhaseUseCase struct {
	repo repo.GamePhaseRepo
	log  *log.Helper
}

func NewGamePhaseUseCase(r repo.GamePhaseRepo, logger log.Logger) *GamePhaseUseCase {
	return &GamePhaseUseCase{
		repo: r,
		log:  log.NewHelper(log.With(logger, "module", "use/game_phase")),
	}
}

func (uc *GamePhaseUseCase) Subscribe() chan string {
	return uc.repo.Subscribe()
}

func (uc *GamePhaseUseCase) Unsubscribe(ch chan string) {
	uc.repo.Unsubscribe(ch)
}
