package user

import (
	"context"
	"go-chat/internal/http/response"
)

type UserService interface {
	GetById(ctx context.Context, userId uint64) (*response.UserResponse, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) GetById(ctx context.Context, userId uint64) (*response.UserResponse, error) {
	usr, err := us.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, response.NewNotFoundErr("user not found", err)
	}

	return &response.UserResponse{
		ID:     userId,
		Name:   usr.Username,
		ImgUrl: usr.ImgUrl,
	}, nil
}
