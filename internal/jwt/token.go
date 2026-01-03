package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(jwtClaims jwt.MapClaims, secret string) (string, error) {
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwtClaims,
	)
	return t.SignedString([]byte(secret))
}

func ValidateJWT(currToken string, secret string) (jwt.MapClaims, error) {
	j, err := jwt.Parse(currToken, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, _ := j.Claims.(jwt.MapClaims)

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func CheckClaims(claims jwt.MapClaims, keys ...string) bool {
	for _, key := range keys {
		if _, ok := claims[key]; !ok {
			return false
		}
	}

	return true
}
