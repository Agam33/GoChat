package auth

import (
	"context"
	"errors"
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
	"go-chat/internal/jwt"
	"go-chat/internal/model"
	"go-chat/internal/utils"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	RefreshToken(string) (response.SignInResponse, error)
	SignUp(context.Context, *request.SignUpRequest) (response.SignInResponse, error)
	SignIn(context.Context, *request.SignInRequst) (response.SignInResponse, error)
}

type authService struct {
	authRepo   AuthRepository
	jwtService jwt.JwtService
}

func NewAuthService(authRepo AuthRepository, jwtService jwt.JwtService) AuthService {
	return &authService{
		authRepo:   authRepo,
		jwtService: jwtService,
	}
}

func (as *authService) RefreshToken(refreshToken string) (response.SignInResponse, error) {
	claims, apperr := as.jwtService.ValidateRefreshToken(refreshToken)
	if apperr != nil {
		return response.SignInResponse{}, apperr
	}

	accessToken, err := as.jwtService.GenerateAccessToken(uint64(claims.UserId))
	if err != nil {
		return response.SignInResponse{}, response.NewInternalServerErr("error generate access token", err)
	}

	newRefreshToken, err := as.jwtService.GenerateRefreshToken(uint64(claims.UserId))
	if err != nil {
		return response.SignInResponse{}, response.NewInternalServerErr("error generate refresh token", err)
	}

	return response.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil

}

func (as *authService) SignUp(ctx context.Context, req *request.SignUpRequest) (response.SignInResponse, error) {
	userId := time.Now().UnixMicro()
	hashPass, ok := utils.HashPassword(req.Password)
	if !ok {
		return response.SignInResponse{}, response.NewInternalServerErr("error hash passsword", nil)
	}

	if err := as.authRepo.SignUp(ctx, &model.User{
		ID:       uint64(userId),
		Username: req.Username,
		Password: hashPass,
	}); err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return response.SignInResponse{}, response.NewBadRequestErr("user already exists", err)
		} else {
			return response.SignInResponse{}, response.NewInternalServerErr("error signup", err)
		}
	}

	accessToken, err := as.jwtService.GenerateAccessToken(uint64(userId))
	if err != nil {
		return response.SignInResponse{}, response.NewInternalServerErr("error generate access token", err)
	}

	refreshToken, err := as.jwtService.GenerateRefreshToken(uint64(userId))
	if err != nil {
		return response.SignInResponse{}, response.NewInternalServerErr("error generate refresh token", err)
	}

	return response.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *authService) SignIn(ctx context.Context, req *request.SignInRequst) (response.SignInResponse, error) {
	usr, err := as.authRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.SignInResponse{}, response.NewNotFoundErr("user not found", err)
		}
		return response.SignInResponse{}, response.NewInternalServerErr("error find user by username (service)", err)
	}

	if !utils.ValidatePassword(usr.Password, req.Password) {
		return response.SignInResponse{}, response.NewBadRequestErr("wrong password", err)
	}

	accessToken, err := as.jwtService.GenerateAccessToken(uint64(usr.ID))
	if err != nil {
		return response.SignInResponse{}, response.NewInternalServerErr("error generate access token", err)
	}

	refreshToken, err := as.jwtService.GenerateRefreshToken(uint64(usr.ID))
	if err != nil {
		return response.SignInResponse{}, response.NewInternalServerErr("error generate refresh token", err)
	}

	return response.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
