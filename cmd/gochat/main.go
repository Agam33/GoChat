package main

import (
	"context"
	"fmt"
	"go-chat/internal/config"
	"go-chat/internal/constant"
	"go-chat/internal/database"
	"go-chat/internal/env"
	"go-chat/internal/http/handler"
	"go-chat/internal/http/middleware"
	"go-chat/internal/http/router"
	"go-chat/internal/jwt"
	"go-chat/internal/rabbitmq"
	"go-chat/internal/services/auth"
	"go-chat/internal/services/chat"
	chatCsmr "go-chat/internal/services/chat/consumer"
	"go-chat/internal/services/room"
	"go-chat/internal/services/user"
	"go-chat/internal/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go-chat/docs"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title GoChat API
// @version	1.0

// @contact.name agam
// @contact.email riswan.dev@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT in the format: Bearer <token>
func main() {
	env, err := env.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	rootCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.NewAppConfig(env)

	psqlDB, err := database.Connect(&cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}

	mqConn, err := rabbitmq.NewConnection(&cfg.RabbitMQ)
	if err != nil {
		log.Fatal(err)
	}

	publisher := rabbitmq.NewPublisher(mqConn)

	wsHub := websocket.NewHub()
	go wsHub.Run()

	r := gin.New()

	// middleware
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.ErrorHandlingMiddleware(),
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	wsHandler := websocket.NewWSHandler(wsHub, publisher, roomService, chatService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)

	// consumer
	chatConsumerDispatcher := chatCsmr.NewChatConsumerHandler(chatService)
	chatConsumer, err := rabbitmq.NewConsumer(mqConn, chatConsumerDispatcher, constant.MQExchangeChat, constant.MQKindTopic, constant.QNameChat, constant.MQBindKeyChat)
	if err != nil {
		log.Fatalf("failed to init chat consumer: %v", err)
	}
	go chatConsumer.Start(rootCtx)

	//router
	router.NewRouter(r, wsHandler, jwtService, userService, authHandler, userHandler, roomHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", env.App.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-rootCtx.Done()
	log.Println("shutting down...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Printf("Server shutdown: %v\n", err)
	}
	log.Println("Server exiting")
}
