package gameflow

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type GameflowRepo interface {
	GetGameflowPhase() (string, error)
	GetGameflowSession() (*resp.GameflowSession, error)
	Reconnect() error
	GetReadyCheckStatus() (*resp.ReadyCheckStatus, error)
	AcceptReadyCheck() error
}

type gameflowRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func NewGameflowRepo(data *infra.Data, global *conf.Global, logger log.Logger) GameflowRepo {
	return &gameflowRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/gameflow")),
	}
}

func (r *gameflowRepo) GetGameflowPhase() (string, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.GameflowPath, nil)
	if err != nil {
		return "", err
	}

	var phase string
	if err := json.Unmarshal(request, &phase); err != nil {
		return "", err
	}

	return phase, nil
}

func (r *gameflowRepo) GetGameflowSession() (*resp.GameflowSession, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.GameflowSessionPath, nil)
	if err != nil {
		return nil, err
	}

	var session resp.GameflowSession
	if err := json.Unmarshal(request, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *gameflowRepo) Reconnect() error {
	_, err := r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.ReconnectPath, nil)
	return err
}

func (r *gameflowRepo) GetReadyCheckStatus() (*resp.ReadyCheckStatus, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.ReadyCheckPath, nil)
	if err != nil {
		return nil, err
	}

	var status resp.ReadyCheckStatus
	if err := json.Unmarshal(request, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (r *gameflowRepo) AcceptReadyCheck() error {
	_, err := r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.MatchmakingPath, nil)
	if err != nil {
		return fmt.Errorf("failed to accept ready check: %w", err)
	}
	r.log.Info("Ready check accepted successfully")
	return nil
}
