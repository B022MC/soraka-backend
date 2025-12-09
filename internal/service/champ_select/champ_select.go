package champ_select

import (
	"go-utils/utils/ecode"
	"go-utils/utils/response"
	"strconv"

	"github.com/B022MC/soraka-backend/internal/biz/champ_select"
	"github.com/gin-gonic/gin"
)

type ChampSelectService struct {
	uc *champ_select.ChampSelectUseCase
}

func NewChampSelectService(uc *champ_select.ChampSelectUseCase) *ChampSelectService {
	return &ChampSelectService{
		uc: uc,
	}
}

func (s *ChampSelectService) RegisterRouter(rootRouter *gin.RouterGroup) {
	router := rootRouter.Group("/champ-select")
	router.GET("/session", s.GetSession)
	router.POST("/select-champion", s.SelectChampion)
	router.POST("/ban-champion", s.BanChampion)
	router.POST("/accept-trade/:tradeId", s.AcceptTrade)
	router.POST("/accept-swap/:swapId", s.AcceptSwap)
	router.POST("/bench-swap/:championId", s.BenchSwap)
	router.GET("/current-champion", s.GetCurrentChampion)
	router.GET("/skin-carousel", s.GetSkinCarousel)
	router.POST("/select-skin", s.SelectSkin)
	router.POST("/reroll", s.Reroll)
}

// GetSession 获取英雄选择会话
func (s *ChampSelectService) GetSession(ctx *gin.Context) {
	session, err := s.uc.GetSession()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取英雄选择会话失败")
		return
	}
	response.Success(ctx, session)
}

// SelectChampion 选择英雄
func (s *ChampSelectService) SelectChampion(ctx *gin.Context) {
	var req struct {
		ActionId   int64 `json:"actionId" binding:"required"`
		ChampionId int   `json:"championId" binding:"required"`
		Completed  bool  `json:"completed"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.SelectChampion(req.ActionId, req.ChampionId, req.Completed); err != nil {
		response.Fail(ctx, ecode.Failed, "选择英雄失败")
		return
	}

	response.Success(ctx, nil)
}

// BanChampion 禁用英雄
func (s *ChampSelectService) BanChampion(ctx *gin.Context) {
	var req struct {
		ActionId   int64 `json:"actionId" binding:"required"`
		ChampionId int   `json:"championId" binding:"required"`
		Completed  bool  `json:"completed"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.BanChampion(req.ActionId, req.ChampionId, req.Completed); err != nil {
		response.Fail(ctx, ecode.Failed, "禁用英雄失败")
		return
	}

	response.Success(ctx, nil)
}

// AcceptTrade 接受英雄交换
func (s *ChampSelectService) AcceptTrade(ctx *gin.Context) {
	tradeId, err := strconv.ParseInt(ctx.Param("tradeId"), 10, 64)
	if err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.AcceptTrade(tradeId); err != nil {
		response.Fail(ctx, ecode.Failed, "接受交换失败")
		return
	}

	response.Success(ctx, nil)
}

// AcceptSwap 接受楼层交换
func (s *ChampSelectService) AcceptSwap(ctx *gin.Context) {
	swapId, err := strconv.ParseInt(ctx.Param("swapId"), 10, 64)
	if err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.AcceptSwap(swapId); err != nil {
		response.Fail(ctx, ecode.Failed, "接受楼层交换失败")
		return
	}

	response.Success(ctx, nil)
}

// BenchSwap 备战席交换
func (s *ChampSelectService) BenchSwap(ctx *gin.Context) {
	championId, err := strconv.Atoi(ctx.Param("championId"))
	if err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.BenchSwap(championId); err != nil {
		response.Fail(ctx, ecode.Failed, "备战席交换失败")
		return
	}

	response.Success(ctx, nil)
}

// GetCurrentChampion 获取当前选择的英雄
func (s *ChampSelectService) GetCurrentChampion(ctx *gin.Context) {
	championId, err := s.uc.GetCurrentChampion()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取当前英雄失败")
		return
	}
	response.Success(ctx, gin.H{"championId": championId})
}

// GetSkinCarousel 获取皮肤轮播
func (s *ChampSelectService) GetSkinCarousel(ctx *gin.Context) {
	skins, err := s.uc.GetSkinCarousel()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取皮肤列表失败")
		return
	}
	response.Success(ctx, skins)
}

// SelectSkin 选择皮肤
func (s *ChampSelectService) SelectSkin(ctx *gin.Context) {
	var req struct {
		SkinId   int  `json:"skinId" binding:"required"`
		Spell1Id *int `json:"spell1Id"`
		Spell2Id *int `json:"spell2Id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, ecode.CaptchaFailed, "参数错误")
		return
	}

	if err := s.uc.SelectSkin(req.SkinId, req.Spell1Id, req.Spell2Id); err != nil {
		response.Fail(ctx, ecode.Failed, "选择皮肤失败")
		return
	}

	response.Success(ctx, nil)
}

// Reroll 摇骰子
func (s *ChampSelectService) Reroll(ctx *gin.Context) {
	if err := s.uc.Reroll(); err != nil {
		response.Fail(ctx, ecode.Failed, "摇骰子失败")
		return
	}
	response.Success(ctx, nil)
}
