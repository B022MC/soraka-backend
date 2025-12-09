package gameflow

import (
	"github.com/B022MC/soraka-backend/internal/dal/repo/gameflow"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/go-kratos/kratos/v2/log"
)

type GameflowUseCase struct {
	repo gameflow.GameflowRepo
	log  *log.Helper
}

func NewGameflowUseCase(repo gameflow.GameflowRepo, logger log.Logger) *GameflowUseCase {
	return &GameflowUseCase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz/gameflow")),
	}
}

func (uc *GameflowUseCase) GetGameflowPhase() (string, error) {
	return uc.repo.GetGameflowPhase()
}

func (uc *GameflowUseCase) GetGameflowSession() (*resp.GameflowSession, error) {
	return uc.repo.GetGameflowSession()
}

func (uc *GameflowUseCase) Reconnect() error {
	return uc.repo.Reconnect()
}

func (uc *GameflowUseCase) GetReadyCheckStatus() (*resp.ReadyCheckStatus, error) {
	return uc.repo.GetReadyCheckStatus()
}

func (uc *GameflowUseCase) AcceptReadyCheck() error {
	return uc.repo.AcceptReadyCheck()
}
