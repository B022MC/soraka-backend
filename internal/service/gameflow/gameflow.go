package gameflow

import (
	"go-utils/utils/ecode"
	"go-utils/utils/response"

	"github.com/B022MC/soraka-backend/internal/biz/gameflow"
	"github.com/gin-gonic/gin"
)

type GameflowService struct {
	uc *gameflow.GameflowUseCase
}

func NewGameflowService(uc *gameflow.GameflowUseCase) *GameflowService {
	return &GameflowService{
		uc: uc,
	}
}

func (s *GameflowService) RegisterRouter(rootRouter *gin.RouterGroup) {
	router := rootRouter.Group("/gameflow")
	router.GET("/phase", s.GetPhase)
	router.GET("/session", s.GetSession)
	router.POST("/reconnect", s.Reconnect)
	router.GET("/ready-check", s.GetReadyCheckStatus)
	router.POST("/accept-ready-check", s.AcceptReadyCheck)
}

// GetPhase 获取当前游戏流程阶段
func (s *GameflowService) GetPhase(ctx *gin.Context) {
	phase, err := s.uc.GetGameflowPhase()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取游戏阶段失败")
		return
	}
	response.Success(ctx, gin.H{"phase": phase})
}

// GetSession 获取游戏流程会话信息
func (s *GameflowService) GetSession(ctx *gin.Context) {
	session, err := s.uc.GetGameflowSession()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取游戏会话失败")
		return
	}
	response.Success(ctx, session)
}

// Reconnect 重新连接游戏
func (s *GameflowService) Reconnect(ctx *gin.Context) {
	if err := s.uc.Reconnect(); err != nil {
		response.Fail(ctx, ecode.Failed, "重新连接失败")
		return
	}
	response.Success(ctx, nil)
}

// GetReadyCheckStatus 获取准备确认状态
func (s *GameflowService) GetReadyCheckStatus(ctx *gin.Context) {
	status, err := s.uc.GetReadyCheckStatus()
	if err != nil {
		response.Fail(ctx, ecode.Failed, "获取准备确认状态失败")
		return
	}
	response.Success(ctx, status)
}

// AcceptReadyCheck 接受对局
func (s *GameflowService) AcceptReadyCheck(ctx *gin.Context) {
	if err := s.uc.AcceptReadyCheck(); err != nil {
		response.Fail(ctx, ecode.Failed, "接受对局失败")
		return
	}
	response.Success(ctx, nil)
}
