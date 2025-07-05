package lcu

import (
	"io"

	gamePhaseBiz "github.com/B022MC/soraka-backend/internal/biz/game_phase"
	"github.com/gin-gonic/gin"
)

// GamePhaseService streams current game phase changes via SSE.
type GamePhaseService struct {
	uc *gamePhaseBiz.GamePhaseUseCase
}

// NewGamePhaseService creates a new GamePhaseService.
func NewGamePhaseService(uc *gamePhaseBiz.GamePhaseUseCase) *GamePhaseService {
	return &GamePhaseService{uc: uc}
}

// RegisterRouter registers the SSE endpoint.
func (s *GamePhaseService) RegisterRouter(root *gin.RouterGroup) {
	group := root.Group("/lcu/gamePhase")
	group.GET("/stream", s.stream)
}

func (s *GamePhaseService) stream(ctx *gin.Context) {
	ch := s.uc.Subscribe()
	defer s.uc.Unsubscribe(ch)

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	ctx.Stream(func(w io.Writer) bool {
		select {
		case phase, ok := <-ch:
			if !ok {
				return false
			}
			ctx.SSEvent("phase", phase)
			return true
		case <-ctx.Done():
			return false
		}
	})
}
