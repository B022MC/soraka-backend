package rank

import (
	"encoding/json"
	"fmt"
	"github.com/B022MC/soraka-backend/consts"
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/B022MC/soraka-backend/internal/infra"
	"github.com/go-kratos/kratos/v2/log"
	"net/http"
)

type RankRepo interface {
	GetRankByPuuid(puuid string) (*resp.Rank, error)
}
type rankRepo struct {
	global *conf.Global
	data   *infra.Data
	log    *log.Helper
}

func (r *rankRepo) GetRankByPuuid(puuid string) (*resp.Rank, error) {
	var rank resp.Rank
	request, err := r.data.LCU.DoRequest(
		http.MethodGet,
		fmt.Sprintf("%s/%s", r.global.Lcu.RankedStatsPath, puuid),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(request, &rank); err != nil {
		return nil, err
	}
	//进行映射中文
	rank.QueueMap.RankedFlexSr.TierCn = consts.TierEnToCn[rank.QueueMap.RankedFlexSr.Tier]
	rank.QueueMap.RankedSolo5x5.TierCn = consts.TierEnToCn[rank.QueueMap.RankedSolo5x5.Tier]
	rank.QueueMap.RankedFlexSr.HighestTierCn = consts.TierEnToCn[rank.QueueMap.RankedFlexSr.Tier]
	rank.QueueMap.RankedSolo5x5.HighestTierCn = consts.TierEnToCn[rank.QueueMap.RankedSolo5x5.Tier]
	rank.QueueMap.RankedFlexSr.QueueTypeCn = consts.QueueTypeToCn[rank.QueueMap.RankedFlexSr.QueueType]
	rank.QueueMap.RankedSolo5x5.QueueType = consts.QueueTypeToCn[rank.QueueMap.RankedSolo5x5.QueueType]
	return &rank, nil
}

func NewRankRepo(data *infra.Data, global *conf.Global, logger log.Logger) RankRepo {
	return &rankRepo{
		data:   data,
		global: global,
		log:    log.NewHelper(log.With(logger, "module", "repo/rank")),
	}
}
