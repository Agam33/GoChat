package middleware

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.Error(response.NewUnauthorized())
			return
		}

		s := strings.Split(bearer, " ")
		if len(s) < 2 || strings.ToLower(s[0]) != "bearer" {
			c.Error(response.NewUnauthorized())
			return
		}

		usr, err := jwtService.ValidateAccessToken(s[1])
		if err != nil {
			c.Error(response.NewUnauthorized())
			return
		}

		c.Set(constant.CtxUserIDKey, usr.UserId)
		c.Next()
	}
}
