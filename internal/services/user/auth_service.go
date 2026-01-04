package user

import (
	"context"
	"go-chat/internal/http/response"
)

type UserService interface {
	GetById(ctx context.Context, userId uint64) (*response.UserReponse, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) GetById(ctx context.Context, userId uint64) (*response.UserReponse, error) {
	usr, err := us.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &response.UserReponse{
		ID:     userId,
		Name:   usr.Name,
		ImgUrl: usr.ImgUrl,
	}, nil
}
