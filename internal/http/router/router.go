package router

import (
	"go-chat/internal/http/handler"
	"go-chat/internal/http/middleware"
	"go-chat/internal/jwt"
	"go-chat/internal/websocket"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	router *gin.Engine,
	wsHandler websocket.WsHandler,
	jwtService jwt.JwtService,
	authHandler handler.AuthHandler,
	userHandler handler.UserHandler,
	roomHandler handler.RoomHandler,
) {
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
		v1 := api.Group("/v1", middleware.JwtMiddleware(jwtService))
		v1.POST("/auth/logout", authHandler.Logout)

		v1.GET("/user/profile", userHandler.GetProfile)
		v1.GET("/user/rooms", userHandler.GetUserRooms)

		v1.GET("/room/create", roomHandler.CreateRoom)
		v1.GET("/room/:id/messages", roomHandler.GetMessages)
		v1.DELETE("/room/:id", roomHandler.DeleteRoom)
		v1.POST("/room/:id/join", roomHandler.JoinRoom)
	}

	router.NoRoute(middleware.NoRouteMiddleware())
}
