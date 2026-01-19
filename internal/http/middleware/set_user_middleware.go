package middleware

import (
	"errors"
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUserMiddleware(userService user.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(constant.CtxUserIDKey)
		if !exists {
			c.AbortWithError(http.StatusUnauthorized, response.NewUnauthorized())
			return
		}

		userId, ok := val.(uint64)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, response.NewInternalServerErr("error: invalid user id type", errors.New("invalid user id type")))
			return
		}

		usr, err := userService.GetById(c.Request.Context(), userId)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Set(constant.CtxUser, usr)
		c.Next()
	}
}
