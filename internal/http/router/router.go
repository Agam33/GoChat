package router

import (
	"go-chat/internal/http/handler"
	"go-chat/internal/http/middleware"
	"go-chat/internal/jwt"
	"go-chat/internal/websocket"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, wsHandler websocket.WsHandler, jwtService jwt.JwtService, authHandler handler.AuthHandler) {
	api := router.Group("/api")

	// websocket
	api.GET("/ws", middleware.JwtMiddleware(jwtService), wsHandler.ServeWS)

	{
		// no authorization
		v1 := api.Group("/v1")
		v1.POST("/auth/signin", authHandler.SignIn)
		v1.POST("/auth/signup", authHandler.SignUp)
		v1.GET("/auth/refresh-token", authHandler.RefreshToken)
	}

	{
		// need authorization
		v1 := api.Group("/v1")
		v1.POST("/auth/logout", authHandler.Logout)
	}

	router.NoRoute(middleware.NoRouteMiddleware())
}
