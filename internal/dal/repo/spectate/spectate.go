package spectate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
)

type SpectateRepo interface {
	SpectateByName(summonerName string) error
}

type spectateRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func NewSpectateRepo(data *infra.Data, global *conf.Global, logger log.Logger) SpectateRepo {
	return &spectateRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/spectate")),
	}
}

func (r *spectateRepo) SpectateByName(summonerName string) error {
	// 首先通过名字获取召唤师信息
	summonerUri := fmt.Sprintf("%s?name=%s", r.global.Lcu.SummonerByNamePath, summonerName)
	summonerData, err := r.data.LCU.DoRequest(http.MethodGet, summonerUri, nil)
	if err != nil {
		return fmt.Errorf("failed to get summoner info: %w", err)
	}

	var summoner map[string]interface{}
	if err := json.Unmarshal(summonerData, &summoner); err != nil {
		return err
	}

	puuid, ok := summoner["puuid"].(string)
	if !ok || puuid == "" {
		return fmt.Errorf("summoner not found")
	}

	// 构造观战请求
	data := map[string]interface{}{
		"allowObserveMode":     "ALL",
		"dropInSpectateGameId": summonerName,
		"gameQueueType":        "",
		"puuid":                puuid,
	}

	body, _ := json.Marshal(data)
	result, err := r.data.LCU.DoRequest(http.MethodPost, r.global.Lcu.SpectateLaunchPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to spectate: %w", err)
	}

	// 如果返回非空响应，表示召唤师不在游戏中
	if len(result) > 0 {
		return fmt.Errorf("summoner is not in game")
	}

	r.log.Infof("Started spectating summoner: %s", summonerName)
	return nil
}
