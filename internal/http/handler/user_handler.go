package handler

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
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

func (h *userHandler) GetProfile(c *gin.Context) {
	userId := c.GetInt64(constant.CtxUserIDKey)
	if userId == 0 {
		c.Error(response.NewUnauthorized())
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
