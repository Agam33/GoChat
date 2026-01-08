package utils

import (
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
