package router

import (
	"go-chat/internal/http/handler"
	"go-chat/internal/http/middleware"
	"go-chat/internal/jwt"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, wsHandler handler.WsHandler, jwtService jwt.JwtService) {
	api := router.Group("/api")

	{
		v1 := api.Group("/v1")

		// websocket
		v1.GET("/ws/chat", middleware.JwtMiddleware(jwtService), wsHandler.ServeChatWS)
	}
}
