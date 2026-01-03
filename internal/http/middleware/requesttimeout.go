package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTimeOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Minute*2)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
