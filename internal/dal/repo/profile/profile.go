package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type ProfileRepo interface {
	SetProfileBackground(skinId int) error
	SetOnlineStatus(message string) error
	SetTierShowed(queue, tier, division string) error
	SetOnlineAvailability(availability string) error
	RemoveTokens() error
	RemovePrestigeCrest() error
	SetProfileIcon(iconId int) error
}

type profileRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func NewProfileRepo(data *infra.Data, global *conf.Global, logger log.Logger) ProfileRepo {
	return &profileRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/profile")),
	}
}

func (r *profileRepo) SetProfileBackground(skinId int) error {
	data := map[string]interface{}{
		"key":   "backgroundSkinId",
		"value": skinId,
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.SummonerProfilePath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set profile background: %w", err)
	}

	r.log.Infof("Profile background set to skin %d", skinId)
	return nil
}

func (r *profileRepo) SetOnlineStatus(message string) error {
	data := map[string]interface{}{
		"statusMessage": message,
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPut, r.global.Lcu.ChatMePath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set online status: %w", err)
	}

	r.log.Infof("Online status set to: %s", message)
	return nil
}

func (r *profileRepo) SetTierShowed(queue, tier, division string) error {
	data := map[string]interface{}{
		"lol": map[string]interface{}{
			"rankedLeagueQueue":    queue,
			"rankedLeagueTier":     tier,
			"rankedLeagueDivision": division,
		},
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPut, r.global.Lcu.ChatMePath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set tier showed: %w", err)
	}

	r.log.Infof("Tier showed set to: %s %s %s", queue, tier, division)
	return nil
}

func (r *profileRepo) SetOnlineAvailability(availability string) error {
	data := map[string]interface{}{
		"availability": availability,
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPut, r.global.Lcu.ChatMePath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set online availability: %w", err)
	}

	r.log.Infof("Online availability set to: %s", availability)
	return nil
}

func (r *profileRepo) RemoveTokens() error {
	// 首先获取当前勋章信息
	chatMe, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.ChatMePath, nil)
	if err != nil {
		return fmt.Errorf("failed to get chat me: %w", err)
	}

	var chatData map[string]interface{}
	if err := json.Unmarshal(chatMe, &chatData); err != nil {
		return err
	}

	// 获取当前选择的勋章
	lolData, ok := chatData["lol"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid chat data structure")
	}

	bannerIdSelected := ""
	if banner, ok := lolData["bannerIdSelected"].(string); ok {
		bannerIdSelected = banner
	}

	data := map[string]interface{}{
		"challengeIds": []int{},
		"bannerAccent": bannerIdSelected,
	}

	body, _ := json.Marshal(data)
	_, err = r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.ChallengesPreferencesPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to remove tokens: %w", err)
	}

	r.log.Info("Challenge tokens removed")
	return nil
}

func (r *profileRepo) RemovePrestigeCrest() error {
	// 获取当前regalia信息
	regaliaData, err := r.data.LCU.DoRequest(http.MethodGet, r.global.Lcu.RegaliaPath, nil)
	if err != nil {
		return fmt.Errorf("failed to get regalia: %w", err)
	}

	var regalia map[string]interface{}
	if err := json.Unmarshal(regaliaData, &regalia); err != nil {
		return err
	}

	bannerType := "lastSeasonHighestRank"
	if bt, ok := regalia["preferredBannerType"].(string); ok {
		bannerType = bt
	}

	data := map[string]interface{}{
		"preferredCrestType":    "prestige",
		"preferredBannerType":   bannerType,
		"selectedPrestigeCrest": 22,
	}

	body, _ := json.Marshal(data)
	_, err = r.data.LCU.DoRequest(http.MethodPut, r.global.Lcu.RegaliaPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to remove prestige crest: %w", err)
	}

	r.log.Info("Prestige crest removed")
	return nil
}

func (r *profileRepo) SetProfileIcon(iconId int) error {
	data := map[string]interface{}{
		"profileIconId": iconId,
	}

	body, _ := json.Marshal(data)
	_, err := r.data.LCU.DoRequest(http.MethodPut, r.global.Lcu.SummonerIconPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set profile icon: %w", err)
	}

	r.log.Infof("Profile icon set to %d", iconId)
	return nil
}
