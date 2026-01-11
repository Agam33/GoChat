package middleware

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		if t, err := c.Cookie("refreshToken"); err == nil {
			token = t
		}

		if token == "" {
			auth := c.GetHeader("X-Refresh-Token")
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnauthorized())
			return
		}

		c.Set(constant.CtxRefreshToken, token)
		c.Next()
	}
}
