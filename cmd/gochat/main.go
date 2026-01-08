package main

import (
	"fmt"
	"go-chat/internal/config"
	"go-chat/internal/database"
	"go-chat/internal/env"
	"go-chat/internal/http/handler"
	"go-chat/internal/http/middleware"
	"go-chat/internal/http/router"
	"go-chat/internal/jwt"
	"go-chat/internal/services/auth"
	"go-chat/internal/services/chat"
	"go-chat/internal/services/room"
	"go-chat/internal/services/user"
	"go-chat/internal/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	env, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewAppConfig(env)

	psqlDB, err := database.Connect(&cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	wsHub := websocket.NewHub()
	go wsHub.Run()

	r := gin.New()

	// middleware
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(middleware.ErrorHandlingMiddleware())

	// repos
	authRepo := auth.NewAuthReposeitory(psqlDB)
	roomRepo := room.NewRoomRepository(psqlDB)
	userRepo := user.NewUserRepository(psqlDB)
	chatRepo := chat.NewChatRepository(psqlDB)

	// services
	jwtService := jwt.NewJwtService(&cfg.JWT)
	authService := auth.NewAuthService(authRepo, jwtService)
	userService := user.NewUserService(userRepo)
	roomService := room.NewRoomService(roomRepo)
	chatService := chat.NewChatService(chatRepo)

	// handlers
	wsHandler := websocket.NewWSHandler(wsHub, userService, roomService, chatService)
	authHandler := handler.NewAuthHandler(authService)

	//router
	router.NewRouter(r, wsHandler, jwtService, authHandler)

	r.Run(fmt.Sprintf(":%d", env.App.Port))
}
