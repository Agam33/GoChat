package room

import (
	"context"
	"encoding/json"
	"errors"
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
	"go-chat/internal/model"
	"go-chat/internal/utils/types"
	"time"

	"gorm.io/gorm"
)

type RoomService interface {
	DeleteRoom(context.Context, uint64) (response.BoolResponse, error)
	GetMessages(context.Context, int64, *types.Pagination) ([]response.RoomMessageResponse, error)
	CreateRoom(context.Context, int64, *request.CreateRoomRequest) (response.BoolResponse, error)
	GetRoomById(context.Context, int64) (response.GetDetailRoomResponse, error)
}

type roomService struct {
	roomRepo RoomRepository
}

func NewRoomService(roomRepo RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (r *roomService) DeleteRoom(ctx context.Context, roomId uint64) (response.BoolResponse, error) {
	if err := r.roomRepo.DeleteRoom(ctx, uint64(roomId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BoolResponse{}, response.NewNotFoundErr("room not found", err)
		}
		return response.BoolResponse{}, response.NewInternalServerErr(err.Error(), err)
	}

	return response.BoolResponse{
		Data: true,
	}, nil
}

func (r *roomService) GetMessages(ctx context.Context, roomId int64, pagination *types.Pagination) ([]response.RoomMessageResponse, error) {
	messages, err := r.roomRepo.GetRoomMessages(ctx, uint64(roomId), pagination)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []response.RoomMessageResponse{}, nil
		}

		return []response.RoomMessageResponse{}, response.NewInternalServerErr(err.Error(), err)
	}

	res := make([]response.RoomMessageResponse, 0, len(messages))
	for _, message := range messages {
		var sender response.UserResponse
		if message.Sender != nil {
			sender = response.UserResponse{
				ID:     message.Sender.ID,
				Name:   message.Sender.Username,
				ImgUrl: message.Sender.ImgUrl,
			}
		}

		res = append(res, response.RoomMessageResponse{
			ID:          message.ID,
			ContentType: message.Type,
			Sender:      sender,
			Content:     json.RawMessage(message.Content),
			CreatedAt:   message.CreatedAt,
		})
	}

	return res, nil
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

func (r *roomService) GetRoomById(ctx context.Context, roomId int64) (response.GetDetailRoomResponse, error) {
	room, err := r.roomRepo.GetById(ctx, uint64(roomId))
	if err != nil {
		return response.GetDetailRoomResponse{}, response.NewNotFoundErr("room not found", err)
	}

	return response.GetDetailRoomResponse{
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
