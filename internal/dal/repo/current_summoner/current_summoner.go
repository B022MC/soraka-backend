package current_summoner

import (
	"encoding/json"
	"fmt"
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"net/url"
)

type CurrentSummonerRepo interface {
	GetCurrentSummonerByPuuid(puuid string) (*resp.Summoner, error)
	GetCurrentSummoner() (*resp.Summoner, error)
	GetCurrentSummonerByName(name string) (*resp.Summoner, error)
}

type currentSummonerRepo struct {
	global     *conf.Global
	data       *infra.Data
	puuidCache *lru.Cache
	log        *log.Helper
}

func NewCurrentSummonerRepo(data *infra.Data, global *conf.Global, logger log.Logger) CurrentSummonerRepo {
	cache, err := lru.New(100) // 可调整容量
	if err != nil {
		panic(fmt.Sprintf("failed to create puuid cache: %v", err))
	}
	return &currentSummonerRepo{
		data:       data,
		global:     global,
		puuidCache: cache,
		log:        log.NewHelper(log.With(logger, "module", "repo/current_summoner")),
	}
}
func (c *currentSummonerRepo) GetCurrentSummoner() (*resp.Summoner, error) {
	var summoner resp.Summoner
	request, err := c.data.LCU.DoRequest(http.MethodGet, c.global.Lcu.SummonerPath, nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(request, &summoner); err != nil {
		return nil, err
	}
	summoner.Dto()
	return &summoner, nil
}
func (c *currentSummonerRepo) GetCurrentSummonerByName(name string) (*resp.Summoner, error) {
	var summoner resp.Summoner
	request, err := c.data.LCU.DoRequest(
		http.MethodGet,
		fmt.Sprintf("%s?name=%s", c.global.Lcu.SummonerByNamePath, url.QueryEscape(name)),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(request, &summoner); err != nil {
		return nil, err
	}
	summoner.Dto()
	return &summoner, nil
}

func (c *currentSummonerRepo) GetCurrentSummonerByPuuid(puuid string) (*resp.Summoner, error) {
	if cached, ok := c.puuidCache.Get(puuid); ok {
		if summoner, ok := cached.(*resp.Summoner); ok {
			c.log.Infof("puuid=%s 缓存命中", puuid)
			return summoner, nil
		}
	}

	uri := fmt.Sprintf("%s/%s", c.global.Lcu.SummonerPuuidPath, puuid)
	request, err := c.data.LCU.DoRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	var summoner resp.Summoner
	if err := json.Unmarshal(request, &summoner); err != nil {
		return nil, err
	}

	c.puuidCache.Add(puuid, &summoner)
	summoner.Dto()
	return &summoner, nil
}
