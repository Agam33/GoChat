package utils

import (
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/utils/types"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string, old string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(old)) == nil
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

func GetPagination(c *gin.Context) (*types.Pagination, error) {
	limitq := c.Query("limit")
	pageq := c.Query("page")

	limit, err := strconv.Atoi(limitq)
	if err != nil {
		return nil, response.NewBadRequestErr("limit should be int", err)
	}

	page, err := strconv.Atoi(pageq)
	if err != nil {
		return nil, response.NewBadRequestErr("page should be int", err)
	}

	if limit <= 0 {
		limit = 10
	}

	return &types.Pagination{
		Limit: limit,
		Page:  page,
	}, nil
}

func PageOffset(page int, limit int) int {
	return (page - 1) * limit
}

func GetUserID(c *gin.Context) (uint64, error) {
	usrIdAny, isExists := c.Get(constant.CtxUserIDKey)
	if !isExists {
		return 0, response.NewUnauthorized()
	}

	userId, ok := usrIdAny.(uint64)
	if !ok {
		return 0, response.NewUnauthorized()
	}

	return userId, nil
}
