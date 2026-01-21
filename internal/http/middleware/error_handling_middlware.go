package middleware

import (
	"errors"
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		code := http.StatusInternalServerError
		status := constant.StatusError
		msg := "internal server error"

		var apperr *response.AppErr
		if errors.As(err, &apperr) {
			code = apperr.Code
			status = constant.StatusFailed

			if code >= 500 {
				log.Printf("[INTERNAL ERROR]: %v", err)
			}

			if code >= 400 && code < 500 {
				msg = apperr.Message
			}
		}

		c.AbortWithStatusJSON(code, gin.H{
			"code":    code,
			"status":  status,
			"message": msg,
		})
	}
}
