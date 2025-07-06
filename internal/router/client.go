package router

import (
	clientService "github.com/B022MC/soraka-backend/internal/service/client"
	"github.com/gin-gonic/gin"
)

type ClientRouter struct {
	clientInfoService *clientService.ClientInfoService
}

func (r *ClientRouter) InitRouter(root *gin.RouterGroup) {
	r.clientInfoService.RegisterRouter(root)

}

func NewClientRouter(
	clientInfoService *clientService.ClientInfoService,
) *ClientRouter {
	return &ClientRouter{
		clientInfoService: clientInfoService,
	}
}
