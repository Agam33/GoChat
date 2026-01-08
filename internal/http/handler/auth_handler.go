package handler

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
	"go-chat/internal/services/auth"
	"go-chat/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Logout(c *gin.Context)
	RefreshToken(c *gin.Context)
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
}

type authHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (h *authHandler) Logout(c *gin.Context) {
	c.SetCookie(
		"refreshToken",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, response.SuccessReponse[bool]{
		Message: constant.StatusSuccess,
		Data:    true,
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	t, err := c.Cookie("refreshToken")
	if err != nil {
		c.Error(response.NewUnauthorized())
		return
	}

	if t == "" {
		var req request.GetRefreshToken

		if err := c.BindJSON(&req); err != nil {
			c.Error(response.NewBadRequestErr("need request token", err))
			return
		}

		t = req.RefreshToken
	}

	resp, err := h.authService.RefreshToken(t)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SetRefreshTokenCookie(c, resp.RefreshToken)

	c.JSON(http.StatusOK, response.SuccessReponse[response.SignInResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}

func (h *authHandler) SignIn(c *gin.Context) {
	var req request.SignInRequst
	if err := c.BindJSON(&req); err != nil {
		c.Error(response.NewBadRequestErr("invalid request signIn", err))
		return
	}

	resp, err := h.authService.SignIn(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SetRefreshTokenCookie(c, resp.RefreshToken)

	c.JSON(http.StatusOK, response.SuccessReponse[response.SignInResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}

func (h *authHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest
	if err := c.BindJSON(&req); err != nil {
		c.Error(response.NewBadRequestErr("invalid request signUp", err))
		return
	}

	resp, err := h.authService.SignUp(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.SetRefreshTokenCookie(c, resp.RefreshToken)

	c.JSON(http.StatusOK, response.SuccessReponse[response.SignInResponse]{
		Message: constant.StatusSuccess,
		Data:    resp,
	})
}
