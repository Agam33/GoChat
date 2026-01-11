package middleware

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AccessTokenMiddleware(jwtService jwt.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader(constant.Authorization)
		if bearer == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnauthorized())
			return
		}

		s := strings.Split(bearer, " ")
		if len(s) < 2 || strings.ToLower(s[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewUnauthorized())
			return
		}

		usr, err := jwtService.ValidateAccessToken(s[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		c.Set(constant.CtxUserIDKey, usr.UserId)
		c.Next()
	}
}
