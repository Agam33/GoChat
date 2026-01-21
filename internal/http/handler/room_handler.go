package handler

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
	"go-chat/internal/services/room"
	"go-chat/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler interface {
	GetRoom(c *gin.Context)
	JoinRoom(c *gin.Context)
	CreateRoom(c *gin.Context)
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

// @Summary Get Room
// @Security BearerAuth
// @Tags	Room
// @Param 	id path	int	true	"Room ID"
// @Success	200 {object}	response.GetDetailRoomResponse
// @Failure 400 {object}	response.AppErr
// @Router	/room/{id}	[get]
func (h *roomHandler) GetRoom(c *gin.Context) {
	roomIdp := c.Param("id")
	if roomIdp == "" {
		c.Error(response.NewBadRequestErr("missing id param", nil))
		return
	}

	roomId, err := strconv.ParseInt(roomIdp, 10, 0)
	if err != nil {
		c.Error(response.NewBadRequestErr("can't parse room id. Id should be a number", nil))
		return
	}

	resp, err := h.roomService.GetRoomById(c.Request.Context(), roomId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponse[response.GetDetailRoomResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}

func (h *roomHandler) JoinRoom(c *gin.Context) {
	var req request.JoinRoomRequst
	if err := c.ShouldBind(&req); err != nil {
		c.Error(response.NewBadRequestErr("invalid request join room", err))
		return
	}

	userId := c.GetUint64(constant.CtxUserIDKey)

	resp, err := h.roomService.JoinRoom(c.Request.Context(), req.RoomId, userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponse[response.BoolResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}

// @Summary Create Room
// @Security BearerAuth
// @Tags	Room
// @Param	request	body request.CreateRoomRequest true "Create Room Request Payload"
// @Success	200 {object}	response.GetRoomResponse
// @Failure 400 {object}	response.AppErr
// @Router	/room/create	[post]
func (h *roomHandler) CreateRoom(c *gin.Context) {
	var req request.CreateRoomRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Error(response.NewBadRequestErr("invalid request create room ", err))
		return
	}

	userId := c.GetUint64(constant.CtxUserIDKey)

	resp, err := h.roomService.CreateRoom(c.Request.Context(), userId, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponse[response.GetRoomResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}

// @Summary Delete Room
// @Security BearerAuth
// @Tags	Room
// @Param	id path int true "Room ID"
// @Success	200 {object}	response.BoolResponse
// @Failure 400 {object}	response.AppErr
// @Router	/room/{id} [post]
func (h *roomHandler) DeleteRoom(c *gin.Context) {
	roomIdq := c.Param("id")
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

// @Summary Get Room Messages
// @Security BearerAuth
// @Tags	Room
// @Param	id path int true "Room ID"
// @Success	200 {array}	response.RoomMessageResponse
// @Failure 400 {object}	response.AppErr
// @Router	/room/{id}/messages [get]
func (h *roomHandler) GetMessages(c *gin.Context) {
	roomIdq := c.Param("id")
	if roomIdq == "" {
		c.Error(response.NewBadRequestErr("room id not found", nil))
		return
	}

	roomId, err := strconv.Atoi(roomIdq)
	if err != nil {
		c.Error(response.NewBadRequestErr("room id should int", err))
		return
	}

	pagination, meta, err := utils.GetPagination(c)
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
		Meta:    meta,
	})
}
