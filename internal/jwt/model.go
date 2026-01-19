package jwt

import "github.com/golang-jwt/jwt/v5"

type Jwtuser struct {
	UserId uint64
	Claims *jwt.MapClaims
}
