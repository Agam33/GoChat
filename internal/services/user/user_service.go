package user

import (
	"context"
	"errors"
	"go-chat/internal/http/response"
	"go-chat/internal/utils/types"

	"gorm.io/gorm"
)

type UserService interface {
	GetUserRooms(context.Context, uint64, *types.Pagination) ([]response.GetRoomResponse, error)
	GetById(ctx context.Context, userId uint64) (response.UserResponse, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (us *userService) GetUserRooms(ctx context.Context, userId uint64, pagination *types.Pagination) ([]response.GetRoomResponse, error) {
	rooms, err := us.userRepo.GetUserRooms(ctx, userId, pagination)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.GetRoomResponse{}, nil
		}
		return []response.GetRoomResponse{}, response.NewInternalServerErr(err.Error(), err)
	}

	res := make([]response.GetRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		res = append(res, response.GetRoomResponse{
			ID:        room.Room.ID,
			Name:      room.Room.Name,
			ImgUrl:    room.Room.ImgUrl,
			CreatedAt: room.CreatedAt,
		})
	}

	return res, nil
}

func (us *userService) GetById(ctx context.Context, userId uint64) (response.UserResponse, error) {
	usr, err := us.userRepo.GetById(ctx, userId)
	if err != nil {
		return response.UserResponse{}, response.NewNotFoundErr("user not found", err)
	}

	return response.UserResponse{
		ID:     userId,
		Name:   usr.Username,
		ImgUrl: usr.ImgUrl,
	}, nil
}
