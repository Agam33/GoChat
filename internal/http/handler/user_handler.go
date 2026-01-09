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

func (h *userHandler) GetUserRooms(c *gin.Context) {
	userId, err := utils.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	pagination, err := utils.GetPagination(c)
	if err != nil {
		c.Error(err)
		return
	}

	resp, err := h.userService.GetUserRooms(c.Request.Context(), uint64(userId), pagination)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponseWithMeta[[]response.GetRoomResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
		Meta: map[string]any{
			"previousPage": pagination.Page - 1,
			"nextPage":     pagination.Page + 1,
			"currPage":     pagination.Page,
		},
	})
}

func (h *userHandler) GetProfile(c *gin.Context) {
	userId, err := utils.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	resp, err := h.userService.GetById(c.Request.Context(), uint64(userId))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.SuccessReponse[response.UserResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}
