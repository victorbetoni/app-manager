package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/victorbetoni/justore/app-manager/internal/config"
	"github.com/victorbetoni/justore/app-manager/internal/infra/http/handler"
	"github.com/victorbetoni/justore/app-manager/internal/infra/jwt"
	"github.com/victorbetoni/justore/app-manager/internal/infra/ws"
)

func Build(connectionHub *ws.ConnectionHub) *gin.Engine {

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, v := range config.GetConfig().Origins {
				if v == origin {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	engine.GET("/ws", UpgraderHandler(connectionHub))
	engine.GET("/app/grid", Authorize(handler.ListApps))
	engine.GET("/app/:appId/state", Authorize(handler.ContainerState))
	engine.POST("/app/:appId/stop", Authorize(handler.ContainerStop))
	engine.POST("/app/:appId/restart", Authorize(handler.ContainerRestart))
	engine.POST("/app/:appId/start", Authorize(handler.ContainerStart))
	engine.POST("/app/:appId/sync", Authorize(handler.SynchronizeProject))

	return engine
}

func Authorize(f gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {

		if !config.GetConfig().UseAuth {
			f(c)
			return
		}

		token, err := c.Cookie(config.GetConfig().Jwt.CookieKey)

		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		_, isADM, err := jwt.ExtractUserIdentifier(token, c.ClientIP(), c.GetHeader("User-Agent"))
		if err != nil || isADM == 0 {
			c.AbortWithStatus(401)
			return
		}

		f(c)

	}
}
