package jwt

import (
	"go-chat/internal/http/response"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateAccessToken(userId uint64) (string, error)
	GenerateRefreshToken(userId uint64) (string, error)
	ValidateRefreshToken(token string) (*Jwtuser, error)
	ValidateAccessToken(token string) (*Jwtuser, error)
}

type jwtService struct {
	jwtConfig *JwtConfig
}

func NewJwtService(jwtConfig *JwtConfig) JwtService {
	return &jwtService{
		jwtConfig: jwtConfig,
	}
}

func (s *jwtService) GenerateAccessToken(userId uint64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(s.jwtConfig.AccessExpire).Unix(),
	}

	return GenerateJWT(claims, s.jwtConfig.AccessSecret)
}

func (s *jwtService) GenerateRefreshToken(userId uint64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(s.jwtConfig.RefreshExpire).Unix(),
	}

	return GenerateJWT(claims, s.jwtConfig.RefreshSecret)
}

func (s *jwtService) ValidateRefreshToken(token string) (*Jwtuser, error) {
	claims, err := ValidateJWT(token, s.jwtConfig.RefreshSecret)
	if err != nil {
		return nil, response.NewUnauthorized()
	}

	if !CheckClaims(claims, "userId", "exp") {
		return nil, response.NewBadRequestErr("invalid token", err)
	}

	uidFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, response.NewBadRequestErr("invalid userId claim", nil)
	}

	return &Jwtuser{
		UserId: int64(uidFloat),
		Claims: &claims,
	}, nil
}

func (s *jwtService) ValidateAccessToken(token string) (*Jwtuser, error) {
	claims, err := ValidateJWT(token, s.jwtConfig.AccessSecret)
	if err != nil {
		return nil, response.NewUnauthorized()
	}

	if !CheckClaims(claims, "userId", "exp") {
		return nil, response.NewBadRequestErr("invalid token", err)
	}

	uidFloat, ok := claims["userId"].(float64)
	if !ok {
		return nil, response.NewBadRequestErr("invalid userId claim", nil)
	}

	return &Jwtuser{
		UserId: int64(uidFloat),
		Claims: &claims,
	}, nil
}
