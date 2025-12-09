package lcu

import (
	"go-utils/utils/ecode"
	"go-utils/utils/response"

	req2 "github.com/B022MC/soraka-backend/internal/dal/req"
	_ "github.com/B022MC/soraka-backend/internal/dal/resp"
	"github.com/gin-gonic/gin"

	currentSummonerUseCaseBiz "github.com/B022MC/soraka-backend/internal/biz/current_summoner"
)

type CurrentSummonerService struct {
	uc *currentSummonerUseCaseBiz.CurrentSummonerUseCase
}

func NewCurrentSummonerService(uc *currentSummonerUseCaseBiz.CurrentSummonerUseCase) *CurrentSummonerService {
	return &CurrentSummonerService{
		uc: uc,
	}
}
func (s *CurrentSummonerService) RegisterRouter(rootRouter *gin.RouterGroup) {
	privateRouter := rootRouter.Group("/lcu")
	privateRouter.GET("/current-summoner", s.GetCurrentSummoner)
	privateRouter.GET("/getRankInfo", s.GetRankInfo)
}

// GetRankInfo
// @Summary 获取召唤师评分信息
// @Description 根据 name 和 puuid 获取召唤师基础信息及其段位信息
// @Tags lcu/GetRankInfo
// @Accept json
// @Produce json
// @Param name query string true "召唤师名（name）"
// @Param puuid query string true "标签（puuid）"
// @Success 200 {object} response.Body{data=resp.SummonerInfo,msg=string}
// @Failure 400 {object} response.Body{msg=string}
// @Failure 500 {object} response.Body{msg=string}
// @Router /lcu/getRankInfo [get]
func (s *CurrentSummonerService) GetRankInfo(ctx *gin.Context) {
	req := req2.SummonerReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}
	rankInfo, err := s.uc.GetSummonerAndRank(req)
	if err != nil {
		response.Fail(ctx, ecode.Failed, "系统错误")
		return
	}
	response.Success(ctx, rankInfo)
}

// GetCurrentSummoner
// @Summary 获取当前召唤师信息
// @Description 获取当前登录的召唤师信息，包括头像、等级、名称等
// @Tags lcu/CurrentSummoner
// @Accept json
// @Produce json
// @Success 200 {object} response.Body{data=resp.Summoner,msg=string}
// @Failure 500 {object} response.Body{msg=string}
// @Router /lcu/current-summoner [get]
func (s *CurrentSummonerService) GetCurrentSummoner(ctx *gin.Context) {
	summoner, err := s.uc.GetCurrentSummoner()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取召唤师信息失败")
		return
	}
	response.Success(ctx, summoner)
}
