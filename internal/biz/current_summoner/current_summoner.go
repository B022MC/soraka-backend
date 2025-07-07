package current_summoner

import (
	"github.com/B022MC/soraka-backend/consts"
	currentSummonerRepo "github.com/B022MC/soraka-backend/internal/dal/repo/current_summoner"
	"github.com/B022MC/soraka-backend/internal/dal/repo/match"
	"github.com/B022MC/soraka-backend/internal/dal/repo/rank"
	"github.com/B022MC/soraka-backend/internal/dal/req"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/go-kratos/kratos/v2/log"
)

type CurrentSummonerUseCase struct {
	repo             currentSummonerRepo.CurrentSummonerRepo
	rankRepo         rank.RankRepo
	matchHistoryRepo match.MatchHistoryRepo
	gameDetailRepo   match.GameDetailRepo
	log              *log.Helper
}

func NewCurrentSummonerUseCase(
	repo currentSummonerRepo.CurrentSummonerRepo,
	rank rank.RankRepo,
	matchHistoryRepo match.MatchHistoryRepo,
	gameDetailRepo match.GameDetailRepo,
	logger log.Logger,
) *CurrentSummonerUseCase {
	return &CurrentSummonerUseCase{
		repo:             repo,
		rankRepo:         rank,
		matchHistoryRepo: matchHistoryRepo,
		gameDetailRepo:   gameDetailRepo,
		log:              log.NewHelper(log.With(logger, "module", "uc/current_summoner"))}
}

func (uc *CurrentSummonerUseCase) GetSummonerAndRank(param req.SummonerReq) (*resp.SummonerInfo, error) {
	var sumoner resp.Summoner
	if param.Puuid != "" {
		res, err := uc.repo.GetCurrentSummonerByPuuid(param.Puuid)
		if err != nil {
			return nil, err
		}
		sumoner = *res
	} else if param.Name != "" {
		res, err := uc.repo.GetCurrentSummonerByName(param.Name)
		if err != nil {
			return nil, err
		}
		sumoner = *res
	} else {
		res, err := uc.repo.GetCurrentSummoner()
		if err != nil {
			return nil, err
		}
		sumoner = *res
	}
	rankRes, err := uc.rankRepo.GetRankByPuuid(sumoner.Puuid)
	if err != nil {
		return nil, err
	}
	matchHistoryReq := req.MatchHistoryReq{
		Puuid:    sumoner.Puuid,
		BegIndex: 0,
		EndIndex: 49,
	}
	matchHistory, err := uc.matchHistoryRepo.GetMatchHistoryByPuuid(matchHistoryReq)
	if err != nil {
		return nil, err
	}
	platformId := ""
	if len(matchHistory.Games.Games) > 0 {
		platformId = matchHistory.Games.Games[0].PlatformId
	}
	sumoner.PlatformIdCn = consts.SGPServerIdToName[platformId]

	return &resp.SummonerInfo{
		Summoner: sumoner,
		Rank:     *rankRes,
	}, nil
}
