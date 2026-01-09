package handler

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/services/room"
	"go-chat/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler interface {
	DeleteRoom(c *gin.Context)
	GetMessages(c *gin.Context)
}

type roomHandler struct {
	roomService room.RoomService
}

func NewRoomHandler(roomService room.RoomService) RoomHandler {
	return &roomHandler{
		roomService: roomService,
	}
}

func (h *roomHandler) DeleteRoom(c *gin.Context) {
	roomIdq := c.Query("id")
	if roomIdq == "" {
		c.Error(response.NewBadRequestErr("room id not found", nil))
		return
	}

	roomId, err := strconv.ParseUint(roomIdq, 10, 0)
	if err != nil {
		c.Error(response.NewBadRequestErr("room id should int", err))
		return
	}

	resp, err := h.roomService.DeleteRoom(c.Request.Context(), roomId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponse[response.BoolResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}

func (h *roomHandler) CreateRoom(c *gin.Context) {

}

func (h *roomHandler) GetMessages(c *gin.Context) {
	roomIdq := c.Query("id")
	if roomIdq == "" {
		c.Error(response.NewBadRequestErr("room id not found", nil))
		return
	}

	roomId, err := strconv.Atoi(roomIdq)
	if err != nil {
		c.Error(response.NewBadRequestErr("room id should int", err))
		return
	}

	pagination, err := utils.GetPagination(c)
	if err != nil {
		c.Error(err)
		return
	}

	resp, err := h.roomService.GetMessages(c.Request.Context(), int64(roomId), pagination)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponseWithMeta[[]response.RoomMessageResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
		Meta: map[string]any{
			"nextPage":     pagination.Page + 1,
			"previousPage": pagination.Page - 1,
			"currentPage":  pagination.Page,
		},
	})
}
