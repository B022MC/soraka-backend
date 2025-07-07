package match

import (
	"encoding/json"
	"fmt"
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
)

type GameDetailRepo interface {
}
type gameDetailRepo struct {
	global          *conf.Global
	data            *infra.Data
	log             *log.Helper
	gameDetailCache *lru.Cache
}

func NewGameDetailRepo(data *infra.Data, global *conf.Global, logger log.Logger) GameDetailRepo {
	gameDetailCache, err := lru.New(100)
	if err != nil {
		panic(fmt.Sprintf("failed to create gameDetailCache: %v", err))
	}
	return &gameDetailRepo{
		data:            data,
		gameDetailCache: gameDetailCache,
		global:          global,
		log:             log.NewHelper(log.With(logger, "module", "repo/gameDetailCache")),
	}
}
func (g *gameDetailRepo) GetGameDetail(gameId int) (*resp.GameDetail, error) {
	// 尝试从缓存获取
	if cached, ok := g.gameDetailCache.Get(gameId); ok {
		if detail, ok := cached.(resp.GameDetail); ok {
			log.Info("GetGameDetail 缓存命中 gameId: %v", gameId)
			return &detail, nil
		}
	}
	// 缓存未命中，从接口获取
	uri := "/lol-match-history/v1/games/%d"
	var gameDetail resp.GameDetail
	request, err := g.data.LCU.DoRequest(http.MethodGet, fmt.Sprintf(uri, gameId), nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(request, &gameDetail); err != nil {
		return nil, err
	}

	// 存入缓存（仅在成功时缓存）
	g.gameDetailCache.Add(gameId, gameDetail)

	return &gameDetail, nil
}
