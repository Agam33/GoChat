package main

import (
	"fmt"
	"go-chat/internal/config"
	"go-chat/internal/database"
	"go-chat/internal/env"
	"go-chat/internal/http/handler"
	"go-chat/internal/http/router"
	"go-chat/internal/http/websocket"
	"go-chat/internal/jwt"
	"go-chat/internal/services/auth"
	"go-chat/internal/services/chat"
	"go-chat/internal/services/room"
	"go-chat/internal/services/user"
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

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	// repos
	authRepo := auth.NewAuthReposeitory(psqlDB)
	roomRepo := room.NewRoomRepository(psqlDB)
	userRepo := user.NewUserRepository(psqlDB)
	chatRepo := chat.NewChatRepository(psqlDB)

	// services
	jwtService := jwt.NewJwtService(&cfg.JWT)
	authService := auth.NewAuthService(authRepo)
	userService := user.NewUserService(userRepo)
	roomService := room.NewRoomService(roomRepo)
	chatService := chat.NewChatService(chatRepo)

	roomHub := websocket.NewRoomHub(chatService, roomService)
	go roomHub.Run()

	// handlers
	wsHandler := handler.NewWSHandler(&userService, roomHub)

	//router
	router.SetupRouter(r, wsHandler, jwtService)

	r.Run(fmt.Sprintf(":%d", env.App.Port))
}
