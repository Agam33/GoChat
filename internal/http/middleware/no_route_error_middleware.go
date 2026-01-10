package middleware

import (
	"go-chat/internal/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRouteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Error(&response.AppErr{
			Code:    http.StatusNotFound,
			Message: "invalid route",
			Err:     nil,
		})
	}
}
