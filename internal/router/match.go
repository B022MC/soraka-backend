package router

import (
	matchService "github.com/B022MC/soraka-backend/internal/service/match"
	"github.com/gin-gonic/gin"
)

type MatchRouter struct {
	matchService *matchService.MatchService
}

func (r *MatchRouter) InitRouter(root *gin.RouterGroup) {
	r.matchService.RegisterRouter(root)
}

func NewMatchRouter(
	matchService *matchService.MatchService,
) *MatchRouter {
	return &MatchRouter{
		matchService: matchService,
	}
}
