package match

import (
	"github.com/B022MC/soraka-backend/internal/dal/repo/match"
	"github.com/B022MC/soraka-backend/internal/dal/req"
	"github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/go-kratos/kratos/v2/log"
)

type MatchUseCase struct {
	matchHistoryRepo match.MatchHistoryRepo
	gameDetailRepo   match.GameDetailRepo
	log              *log.Helper
}

func NewMatchUseCase(
	matchHistoryRepo match.MatchHistoryRepo,
	gameDetailRepo match.GameDetailRepo,
	logger log.Logger,
) *MatchUseCase {
	return &MatchUseCase{
		matchHistoryRepo: matchHistoryRepo,
		gameDetailRepo:   gameDetailRepo,
		log:              log.NewHelper(log.With(logger, "module", "uc/match")),
	}
}

// GetMatchHistory 获取对局历史
func (uc *MatchUseCase) GetMatchHistory(puuid string, begIndex, endIndex int) (*resp.MatchHistory, error) {
	matchHistoryReq := req.MatchHistoryReq{
		Puuid:    puuid,
		BegIndex: begIndex,
		EndIndex: endIndex,
	}

	matchHistory, err := uc.matchHistoryRepo.GetMatchHistoryByPuuid(matchHistoryReq)
	if err != nil {
		uc.log.Errorf("获取对局历史失败: %v", err)
		return nil, err
	}

	return matchHistory, nil
}

// GetGameDetail 获取对局详情
func (uc *MatchUseCase) GetGameDetail(gameId int) (*resp.GameDetail, error) {
	gameDetail, err := uc.gameDetailRepo.GetGameDetail(gameId)
	if err != nil {
		uc.log.Errorf("获取对局详情失败: %v", err)
		return nil, err
	}

	return gameDetail, nil
}
