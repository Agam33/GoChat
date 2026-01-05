package middleware

import (
	"go-chat/internal/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRouteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"status":  constant.StatusError,
			"message": "route not found",
		})
	}
}
