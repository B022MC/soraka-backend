package server

import (
	"github.com/B022MC/soraka-backend/internal/conf"
	"github.com/B022MC/soraka-backend/internal/router"
	"github.com/gin-contrib/cors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/tx7do/kratos-transport/transport/gin"
	"time"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger, router *router.RootRouter) *gin.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := gin.NewServer(
		gin.WithAddress(c.Http.Addr),
		gin.WithLogger(logger),
	)
	initPlugin(srv, logger)
	router.InitRouter(srv.Group(""))
	return srv
}

func initPlugin(srv *gin.Server, logger log.Logger) {
	srv.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
