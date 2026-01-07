package router

import (
	"go-chat/internal/http/middleware"
	"go-chat/internal/jwt"
	"go-chat/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, wsHandler websocket.WsHandler, jwtService jwt.JwtService) {
	api := router.Group("/api")

	// websocket
	api.GET("/ws", middleware.JwtMiddleware(jwtService), wsHandler.ServeWS)
}
