package lcu

import (
	"github.com/gin-gonic/gin"
	"go-utils/utils/response"

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
	privateRouter := rootRouter.Group("/lcu/currentSummoner")
	privateRouter.POST("/getCurrentSummoner", s.GetCurrSummoner)
}
func (s *CurrentSummonerService) GetCurrSummoner(ctx *gin.Context) {
	//data, err := s.uc.GetCurrSummoner()
	//if err != nil {
	//	response.ErrorMsg(ctx, err.Error())
	//	return
	//}
	response.Success(ctx, "获取成功")
}
