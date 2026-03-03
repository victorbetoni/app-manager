package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/victorbetoni/justore/app-manager/internal/config"
	"github.com/victorbetoni/justore/app-manager/internal/infra/jwt"
	"github.com/victorbetoni/justore/app-manager/internal/infra/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func UpgraderHandler(connectionHub *ws.ConnectionHub) func(*gin.Context) {

	return func(c *gin.Context) {

		identifier := uuid.New().String()

		if config.GetConfig().UseAuth {
			cookie, err := c.Cookie(config.GetConfig().Jwt.CookieKey)

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			dcId, adm, err := jwt.ExtractUserIdentifier(cookie, c.ClientIP(), c.GetHeader("User-Agent"))
			if err != nil || adm != 1 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			identifier = dcId
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()

		session := ws.NewSession(ctx, identifier, c.ClientIP(), c.GetHeader("User-Agent"), conn, connectionHub)
		if err := session.Listen(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

	}
}
