package middleware

import (
	"go-chat/internal/constant"
	"go-chat/internal/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			c.AbortWithStatusJSON(401, map[string]any{
				"status":  401,
				"message": "unauthorized",
			})
			return
		}

		s := strings.Split(bearer, " ")
		if len(s) < 2 || strings.ToLower(s[0]) != "bearer" {
			c.AbortWithStatusJSON(401, map[string]any{
				"status":  401,
				"message": "unauthorized",
			})
			return
		}

		usr, err := jwtService.ValidateAccessToken(s[1])
		if err != nil {
			c.AbortWithStatusJSON(401, map[string]any{
				"status":  401,
				"message": "unauthorized",
			})
			return
		}

		c.Set(constant.CtxUserIDKey, usr.UserId)
		c.Next()
	}
}
