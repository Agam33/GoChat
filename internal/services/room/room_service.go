package room

import (
	"context"
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
	"go-chat/internal/model"
	"time"
)

type RoomService interface {
	CreateRoom(context.Context, int64, *request.CreateRoomRequest) (response.BoolResponse, error)
	GetRoomById(context.Context, int64) (response.GetRoomResponse, error)
}

type roomService struct {
	roomRepo RoomRepository
}

func NewRoomService(roomRepo RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (r *roomService) CreateRoom(ctx context.Context, userId int64, req *request.CreateRoomRequest) (response.BoolResponse, error) {
	if req == nil {
		return response.BoolResponse{}, response.NewBadRequestErr("invalid create room request", nil)
	}

	roomId := time.Now().UnixMicro()
	if err := r.roomRepo.CreateRoom(ctx, &model.Room{
		ID:        uint64(roomId),
		CreatorID: uint64(userId),
		Name:      req.Name,
		ImgUrl:    nil,
	}); err != nil {
		return response.BoolResponse{}, err
	}

	return response.BoolResponse{
		Data: true,
	}, nil
}

func (r *roomService) GetRoomById(ctx context.Context, roomId int64) (response.GetRoomResponse, error) {
	room, err := r.roomRepo.GetById(ctx, uint64(roomId))
	if err != nil {
		return response.GetRoomResponse{}, response.NewNotFoundErr("room not found", err)
	}

	return response.GetRoomResponse{
		ID:     room.ID,
		Name:   room.Name,
		ImgUrl: room.ImgUrl,
		Creator: response.UserResponse{
			ID:     room.CreatorID,
			Name:   room.Name,
			ImgUrl: room.ImgUrl,
		},
		CreatedAt: room.CreatedAt,
	}, nil
}
