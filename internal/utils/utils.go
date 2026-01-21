package utils

import (
	"go-chat/internal/utils/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) (string, bool) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err == nil
}

func SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(
		"refreshToken", // name
		refreshToken,   // value
		7*24*3600,      // maxAge (7 days)
		"/",            // path
		"",             // domain (optional)
		false,          // secure (true for https)
		true,           // httpOnly
	)
}

func GetPagination(c *gin.Context) (*types.Pagination, types.Meta, error) {
	limitq := c.Query("limit")
	pageq := c.Query("page")

	limit, err := strconv.Atoi(limitq)
	if err != nil {
		limit = 10
	}

	page, err := strconv.Atoi(pageq)
	if err != nil {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	return &types.Pagination{
			Limit: limit,
			Page:  page,
		}, types.Meta{
			"nextPage": page + 1,
			"prevPage": page - 1,
			"currPage": page,
		}, nil
}

func PageOffset(page int, limit int) int {
	return (page - 1) * limit
}
