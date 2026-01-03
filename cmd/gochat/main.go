package main

import (
	"fmt"
	"go-chat/internal/config"
	"go-chat/internal/database"
	"go-chat/internal/env"
	"go-chat/internal/http/websocket"
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

	roomHub := websocket.NewRoomHub()
	go roomHub.Run()

	router := gin.Default()

	router.Run(fmt.Sprintf(":%d", env.App.Port))
}
