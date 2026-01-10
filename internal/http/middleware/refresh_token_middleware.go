package middleware

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
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
			auth := c.GetHeader(constant.Authorization)
			if strings.HasPrefix(auth, "Bearer ") {
				token = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if token == "" {
			c.Error(response.NewUnauthorized())
			return
		}

		c.Set(constant.CtxRefreshToken, token)
		c.Next()
	}
}
