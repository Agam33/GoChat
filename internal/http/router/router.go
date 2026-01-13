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
	api.GET("/ws", middleware.AccessTokenMiddleware(jwtService), wsHandler.ServeWS)

	{
		// no authorization
		v1 := api.Group("/v1")
		v1.POST("/auth/signin", authHandler.SignIn)
		v1.POST("/auth/signup", authHandler.SignUp)
		v1.GET("/auth/refresh-token", middleware.RefreshTokenMiddleware(), authHandler.RefreshToken)
	}

	{
		// need authorization
		v1 := api.Group("/v1", middleware.AccessTokenMiddleware(jwtService))
		v1.POST("/auth/logout", authHandler.Logout)

		v1.GET("/user/profile", userHandler.GetProfile)
		v1.GET("/user/rooms", userHandler.GetUserRooms)

		v1.POST("/room/create", roomHandler.CreateRoom)
		v1.GET("/room/:id/messages", roomHandler.GetMessages)
		v1.DELETE("/room/:id", roomHandler.DeleteRoom)
		v1.GET("/room/:id", roomHandler.GetRoom)
		v1.POST("/room/:id/join", roomHandler.JoinRoom)
	}

	router.NoRoute(middleware.NoRouteMiddleware())
}
