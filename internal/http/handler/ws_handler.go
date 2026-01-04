package handler

import (
	"fmt"
	"go-chat/internal/constant"
	chws "go-chat/internal/http/websocket"
	"go-chat/internal/services/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler interface {
	ServeChatWS(*gin.Context)
}

type wsHandler struct {
	userService user.UserService
	roomHub     *chws.RoomHub
	upgrader    websocket.Upgrader
}

func NewWSHandler(userService user.UserService, roomHub *chws.RoomHub) WsHandler {
	return &wsHandler{
		userService: userService,
		roomHub:     roomHub,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *wsHandler) ServeChatWS(c *gin.Context) {
	roomIdq := c.Query("room_id")
	if roomIdq == "" {
		c.AbortWithStatusJSON(400, map[string]any{
			"status":  400,
			"message": "invalid query room id",
		})
		return
	}

	roomId, err := strconv.ParseInt(roomIdq, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(500, map[string]any{
			"status":  500,
			"message": "internal server error",
		})
		return
	}

	userId := c.GetUint64(constant.CtxUserIDKey)
	resp, err := h.userService.GetById(c.Request.Context(), userId)
	if err != nil {
		c.AbortWithStatusJSON(400, map[string]any{
			"status":  404,
			"message": "user not found",
		})

		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("can't serve websocket: %v", err)
		c.AbortWithStatusJSON(500, map[string]any{
			"status":  500,
			"message": "internal server error",
		})
		return
	}

	client := &chws.Client{
		UserId: resp.ID,
		Name:   resp.Name,
		ImgUrl: resp.ImgUrl,
		RoomId: roomId,
		Conn:   conn,
		Send:   make(chan []byte, 512),
		Ctx:    c.Request.Context(),
	}

	h.roomHub.Register <- client
}
