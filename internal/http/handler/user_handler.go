package handler

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/services/user"
	"go-chat/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUserRooms(c *gin.Context)
	GetProfile(c *gin.Context)
}

type userHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// @Summary Get User Rooms
// @Security BearerAuth
// @Tags	User
// @Success	200 {object}	response.GetRoomsResponse
// @Failure 400 {object}	response.AppErr
// @Router	/user/rooms [get]
func (h *userHandler) GetUserRooms(c *gin.Context) {
	userId := c.GetUint64(constant.CtxUserIDKey)

	pagination, meta, err := utils.GetPagination(c)
	if err != nil {
		c.Error(err)
		return
	}

	resp, err := h.userService.GetUserRooms(c.Request.Context(), userId, pagination)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponseWithMeta[[]response.GetRoomResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
		Meta:    meta,
	})
}

// @Summary Get Profiles
// @Security BearerAuth
// @Tags	User
// @Success	200 {object}	response.UserResponse
// @Failure 400 {object}	response.AppErr
// @Router	/user/profile [get]
func (h *userHandler) GetProfile(c *gin.Context) {
	userId := c.GetUint64(constant.CtxUserIDKey)

	resp, err := h.userService.GetById(c.Request.Context(), userId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponse[response.UserResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}
