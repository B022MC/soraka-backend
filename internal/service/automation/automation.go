package automation

import (
	"go-utils/utils/ecode"
	"go-utils/utils/response"

	"github.com/B022MC/soraka-backend/internal/biz/automation"
	"github.com/gin-gonic/gin"
)

type AutomationService struct {
	uc *automation.AutomationUseCase
}

func NewAutomationService(uc *automation.AutomationUseCase) *AutomationService {
	return &AutomationService{
		uc: uc,
	}
}

func (s *AutomationService) RegisterRouter(rootRouter *gin.RouterGroup) {
	router := rootRouter.Group("/automation")
	router.POST("/accept-ready-check", s.AutoAcceptReadyCheck)
	router.POST("/select-champion", s.AutoSelectChampion)
	router.POST("/ban-champion", s.AutoBanChampion)
	router.POST("/accept-trades", s.AutoAcceptTrades)
	router.POST("/accept-swaps", s.AutoAcceptSwaps)
	router.POST("/apply-rune-page", s.ApplyRunePage)
}

// AutoAcceptReadyCheck 自动接受对局
func (s *AutomationService) AutoAcceptReadyCheck(ctx *gin.Context) {
	if err := s.uc.AutoAcceptReadyCheck(); err != nil {
		response.Fail(ctx, ecode.Failed, "自动接受对局失败")
		return
	}
	response.Success(ctx, nil)
}

// AutoSelectChampion 自动选择英雄
func (s *AutomationService) AutoSelectChampion(ctx *gin.Context) {
	var req struct {
		ChampionId int `json:"championId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.AutoSelectChampion(req.ChampionId); err != nil {
		response.Fail(ctx, ecode.Failed, "自动选择英雄失败")
		return
	}

	response.Success(ctx, nil)
}

// AutoBanChampion 自动禁用英雄
func (s *AutomationService) AutoBanChampion(ctx *gin.Context) {
	var req struct {
		ChampionId int `json:"championId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.AutoBanChampion(req.ChampionId); err != nil {
		response.Fail(ctx, ecode.Failed, "自动禁用英雄失败")
		return
	}

	response.Success(ctx, nil)
}

// AutoAcceptTrades 自动接受所有英雄交换
func (s *AutomationService) AutoAcceptTrades(ctx *gin.Context) {
	if err := s.uc.AutoAcceptTrades(); err != nil {
		response.Fail(ctx, ecode.Failed, "自动接受交换失败")
		return
	}
	response.Success(ctx, nil)
}

// AutoAcceptSwaps 自动接受所有楼层交换
func (s *AutomationService) AutoAcceptSwaps(ctx *gin.Context) {
	if err := s.uc.AutoAcceptSwaps(); err != nil {
		response.Fail(ctx, ecode.Failed, "自动接受楼层交换失败")
		return
	}
	response.Success(ctx, nil)
}

// ApplyRunePage 应用符文页
func (s *AutomationService) ApplyRunePage(ctx *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required"`
		PrimaryStyleId  int    `json:"primaryStyleId" binding:"required"`
		SubStyleId      int    `json:"subStyleId" binding:"required"`
		SelectedPerkIds []int  `json:"selectedPerkIds" binding:"required,min=6,max=6"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误，符文必须为6个")
		return
	}

	if err := s.uc.ApplyRunePage(req.Name, req.PrimaryStyleId, req.SubStyleId, req.SelectedPerkIds); err != nil {
		response.Fail(ctx, ecode.Failed, "应用符文页失败")
		return
	}

	response.Success(ctx, nil)
}
