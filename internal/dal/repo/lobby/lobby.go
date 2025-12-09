package lobby

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type LobbyRepo interface {
	Create5v5PracticeLobby(lobbyName, password string) error
}

type lobbyRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func NewLobbyRepo(data *infra.Data, global *conf.Global, logger log.Logger) LobbyRepo {
	return &lobbyRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/lobby")),
	}
}

func (r *lobbyRepo) Create5v5PracticeLobby(lobbyName, password string) error {
	data := map[string]interface{}{
		"customGameLobby": map[string]interface{}{
			"configuration": map[string]interface{}{
				"gameMode":         "PRACTICETOOL",
				"gameMutator":      "",
				"gameServerRegion": "",
				"mapId":            11,
				"mutators":         map[string]int{"id": 1},
				"spectatorPolicy":  "AllAllowed",
				"teamSize":         5,
			},
			"lobbyName":     lobbyName,
			"lobbyPassword": password,
		},
		"isCustom": true,
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.LobbyPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create practice lobby: %w", err)
	}

	r.log.Infof("5v5 practice lobby '%s' created", lobbyName)
	return nil
}
