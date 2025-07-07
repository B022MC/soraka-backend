package match

import (
	"encoding/json"
	"fmt"
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/req"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
	lru "github.com/hashicorp/golang-lru"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MatchHistoryRepo interface {
	GetMatchHistoryByPuuid(MatchHistoryReq req.MatchHistoryReq) (*resp.MatchHistory, error)
}
type matchHistoryRepo struct {
	global            *conf.Global
	data              *infra.Data
	log               *log.Helper
	matchHistoryCache *lru.Cache
}
type lruValue struct {
	expiresAt    time.Time
	matchHistory resp.MatchHistory
}

func NewMatchHistoryRepo(data *infra.Data, global *conf.Global, logger log.Logger) MatchHistoryRepo {
	matchHistoryCache, err := lru.New(100)
	if err != nil {
		panic(fmt.Sprintf("failed to create match history cache: %v", err))
	}
	return &matchHistoryRepo{
		data:              data,
		matchHistoryCache: matchHistoryCache,
		global:            global,
		log:               log.NewHelper(log.With(logger, "module", "repo/match_history")),
	}
}
func (m *matchHistoryRepo) GetMatchHistoryByPuuid(MatchHistoryReq req.MatchHistoryReq) (*resp.MatchHistory, error) {
	uri := "/lol-match-history/v1/products/lol/%s/matches?%s"

	// 尝试从缓存获取
	if cached, ok := m.matchHistoryCache.Get(MatchHistoryReq.Puuid); ok {
		if value, ok := cached.(lruValue); ok && value.expiresAt.After(time.Now()) {
			// 检查请求范围是否在缓存范围内 (0-49)
			if MatchHistoryReq.BegIndex >= 0 && MatchHistoryReq.EndIndex <= 49 && MatchHistoryReq.EndIndex < len(value.matchHistory.Games.Games) {
				log.Info("GetMatchHistoryByPuuid() 缓存命中")
				result := value.matchHistory
				result.Games.Games = result.Games.Games[MatchHistoryReq.BegIndex : MatchHistoryReq.EndIndex+1]
				return &result, nil
			}
		}
	}

	// 缓存未命中或范围不匹配，从接口获取
	params := url.Values{}
	if MatchHistoryReq.BegIndex == 0 {
		params.Add("begIndex", "0")
		params.Add("endIndex", "49")
	} else {
		params.Add("begIndex", strconv.Itoa(MatchHistoryReq.BegIndex))
		params.Add("endIndex", strconv.Itoa(MatchHistoryReq.EndIndex))
	}

	var matchHistory resp.MatchHistory
	request, err := m.data.LCU.DoRequest(http.MethodGet, fmt.Sprintf(uri, MatchHistoryReq.Puuid, params.Encode()), nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(request, &matchHistory); err != nil {
		return nil, err
	}

	// 如果获取的是0-49范围的数据，则更新缓存
	if MatchHistoryReq.BegIndex == 0 {
		randomTime := time.Duration(rand.Intn(120)) * time.Second
		value := lruValue{
			expiresAt:    time.Now().Add(time.Minute * 1).Add(randomTime),
			matchHistory: matchHistory,
		}
		m.matchHistoryCache.Add(MatchHistoryReq.Puuid, value)
	}
	if MatchHistoryReq.BegIndex == 0 && MatchHistoryReq.EndIndex < len(matchHistory.Games.Games) {
		matchHistory.Games.Games = matchHistory.Games.Games[MatchHistoryReq.BegIndex : MatchHistoryReq.EndIndex+1]
	}

	return &matchHistory, nil
}
