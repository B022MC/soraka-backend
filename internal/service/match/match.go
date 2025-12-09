package match

import (
	"go-utils/utils/ecode"
	"go-utils/utils/response"
	"strconv"

	matchUseCaseBiz "github.com/B022MC/soraka-backend/internal/biz/match"
	"github.com/gin-gonic/gin"
)

type MatchService struct {
	uc *matchUseCaseBiz.MatchUseCase
}

func NewMatchService(uc *matchUseCaseBiz.MatchUseCase) *MatchService {
	return &MatchService{
		uc: uc,
	}
}

func (s *MatchService) RegisterRouter(rootRouter *gin.RouterGroup) {
	matchRouter := rootRouter.Group("/match")
	matchRouter.GET("/history", s.GetMatchHistory)
	matchRouter.GET("/detail/:gameId", s.GetGameDetail)
}

// GetMatchHistory
// @Summary 获取对局历史
// @Description 获取召唤师的对局历史记录
// @Tags match
// @Accept json
// @Produce json
// @Param puuid query string true "召唤师PUUID"
// @Param beg_index query int false "开始索引" default(0)
// @Param end_index query int false "结束索引" default(19)
// @Success 200 {object} response.Body{data=resp.MatchHistory}
// @Failure 400 {object} response.Body
// @Failure 500 {object} response.Body
// @Router /match/history [get]
func (s *MatchService) GetMatchHistory(ctx *gin.Context) {
	puuid := ctx.Query("puuid")
	if puuid == "" {
		response.Fail(ctx, ecode.CaptchaFailed, "puuid不能为空")
		return
	}

	begIndex, _ := strconv.Atoi(ctx.DefaultQuery("beg_index", "0"))
	endIndex, _ := strconv.Atoi(ctx.DefaultQuery("end_index", "19"))

	matchHistory, err := s.uc.GetMatchHistory(puuid, begIndex, endIndex)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取对局历史失败")
		return
	}

	response.Success(ctx, matchHistory)
}

// GetGameDetail
// @Summary 获取对局详情
// @Description 根据对局ID获取对局详细信息
// @Tags match
// @Accept json
// @Produce json
// @Param gameId path int true "对局ID"
// @Success 200 {object} response.Body{data=resp.GameDetail}
// @Failure 400 {object} response.Body
// @Failure 500 {object} response.Body
// @Router /match/detail/{gameId} [get]
func (s *MatchService) GetGameDetail(ctx *gin.Context) {
	gameIdStr := ctx.Param("gameId")
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "gameId格式错误")
		return
	}

	gameDetail, err := s.uc.GetGameDetail(gameId)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取对局详情失败")
		return
	}

	response.Success(ctx, gameDetail)
}
