package champ_select

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type ChampSelectRepo interface {
	GetSession() (*resp.ChampSelectSession, error)
	SelectChampion(actionId int64, championId int, completed bool) error
	BanChampion(actionId int64, championId int, completed bool) error
	AcceptTrade(tradeId int64) error
	AcceptSwap(swapId int64) error
	BenchSwap(championId int) error
	GetCurrentChampion() (int, error)
	GetSkinCarousel() ([]resp.SkinCarousel, error)
	SelectSkin(skinId int, spell1Id, spell2Id *int) error
	Reroll() error
}

type champSelectRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func NewChampSelectRepo(data *infra.Data, global *conf.Global, logger log.Logger) ChampSelectRepo {
	return &champSelectRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/champ_select")),
	}
}

func (r *champSelectRepo) GetSession() (*resp.ChampSelectSession, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.ChampSelectSessionPath, nil)
	if err != nil {
		return nil, err
	}

	var session resp.ChampSelectSession
	if err := json.Unmarshal(request, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *champSelectRepo) SelectChampion(actionId int64, championId int, completed bool) error {
	data := map[string]interface{}{
		"championId": championId,
		"type":       "pick",
		"completed":  completed,
	}

	body, _ := json.Marshal(data)
	uri := fmt.Sprintf("%s/%d", r.global.Lcu.ChampSelectActionsPath, actionId)

	_, err := r.data.LCU.DoRequest(http.MethodPatch, uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to select champion: %w", err)
	}

	r.log.Infof("Champion %d selected for action %d", championId, actionId)
	return nil
}

func (r *champSelectRepo) BanChampion(actionId int64, championId int, completed bool) error {
	data := map[string]interface{}{
		"championId": championId,
		"type":       "ban",
		"completed":  completed,
	}

	body, _ := json.Marshal(data)
	uri := fmt.Sprintf("%s/%d", r.global.Lcu.ChampSelectActionsPath, actionId)

	_, err := r.data.LCU.DoRequest(http.MethodPatch, uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to ban champion: %w", err)
	}

	r.log.Infof("Champion %d banned for action %d", championId, actionId)
	return nil
}

func (r *champSelectRepo) AcceptTrade(tradeId int64) error {
	uri := fmt.Sprintf("%s/%d/accept", r.global.Lcu.ChampSelectTradesPath, tradeId)
	_, err := r.data.LCU.DoRequest(http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to accept trade: %w", err)
	}

	// Clear the trade
	clearUri := fmt.Sprintf("/lol-champ-select/v1/ongoing-trade/%d/clear", tradeId)
	_, _ = r.data.LCU.DoRequest(http.MethodPost, clearUri, nil)

	r.log.Infof("Trade %d accepted", tradeId)
	return nil
}

func (r *champSelectRepo) AcceptSwap(swapId int64) error {
	uri := fmt.Sprintf("%s/%d/accept", r.global.Lcu.ChampSelectSwapsPath, swapId)
	_, err := r.data.LCU.DoRequest(http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to accept swap: %w", err)
	}

	// Clear the swap
	clearUri := fmt.Sprintf("/lol-champ-select/v1/ongoing-swap/%d/clear", swapId)
	_, _ = r.data.LCU.DoRequest(http.MethodPost, clearUri, nil)

	r.log.Infof("Swap %d accepted", swapId)
	return nil
}

func (r *champSelectRepo) BenchSwap(championId int) error {
	uri := fmt.Sprintf("%s/%d", r.global.Lcu.ChampSelectBenchSwapPath, championId)
	_, err := r.data.LCU.DoRequest(http.MethodPost, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to bench swap: %w", err)
	}

	r.log.Infof("Bench swapped to champion %d", championId)
	return nil
}

func (r *champSelectRepo) GetCurrentChampion() (int, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.CurrentChampionPath, nil)
	if err != nil {
		return 0, err
	}

	var championId int
	if err := json.Unmarshal(request, &championId); err != nil {
		return 0, err
	}

	return championId, nil
}

func (r *champSelectRepo) GetSkinCarousel() ([]resp.SkinCarousel, error) {
	request, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.ChampSelectSkinCarouselPath, nil)
	if err != nil {
		return nil, err
	}

	var skins []resp.SkinCarousel
	if err := json.Unmarshal(request, &skins); err != nil {
		return nil, err
	}

	return skins, nil
}

func (r *champSelectRepo) SelectSkin(skinId int, spell1Id, spell2Id *int) error {
	data := map[string]interface{}{
		"selectedSkinId": skinId,
	}

	if spell1Id != nil {
		data["spell1Id"] = *spell1Id
	}
	if spell2Id != nil {
		data["spell2Id"] = *spell2Id
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPatch, r.global.Lcu.ChampSelectMySelectionPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to select skin: %w", err)
	}

	r.log.Infof("Skin %d selected", skinId)
	return nil
}

func (r *champSelectRepo) Reroll() error {
	_, err := r.data.LCU.DoRequest(http.MethodPost, "/lol-champ-select/v1/session/my-selection/reroll", nil)
	if err != nil {
		return fmt.Errorf("failed to reroll: %w", err)
	}

	r.log.Info("Rerolled successfully")
	return nil
}
