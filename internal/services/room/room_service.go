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
	LeaveRoom(ctx context.Context, roomId uint64, userId uint64) error
	IsJoined(ctx context.Context, roomId uint64, userId uint64) (response.BoolResponse, error)
	JoinRoom(ctx context.Context, roomId uint64, userId uint64) (response.BoolResponse, error)
	DeleteRoom(ctx context.Context, roomId uint64) (response.BoolResponse, error)
	GetMessages(ctx context.Context, roomId int64, pagination *types.Pagination) ([]response.RoomMessageResponse, error)
	CreateRoom(ctx context.Context, userId uint64, req *request.CreateRoomRequest) (response.GetRoomResponse, error)
	GetRoomById(ctx context.Context, roomId int64) (response.GetDetailRoomResponse, error)
}

type roomService struct {
	roomRepo RoomRepository
}

func NewRoomService(roomRepo RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (r *roomService) LeaveRoom(ctx context.Context, roomId uint64, userId uint64) error {
	if err := r.roomRepo.LeaveRoom(ctx, roomId, userId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewNotFoundErr("room not found", err)
		}
		return response.NewInternalServerErr(err.Error(), err)
	}

	return nil
}

func (r *roomService) IsJoined(ctx context.Context, roomdId uint64, userId uint64) (response.BoolResponse, error) {
	if err := r.roomRepo.IsJoined(ctx, roomdId, userId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BoolResponse{Data: true}, response.NewBadRequestErr("room not found", err)
		}

		return response.BoolResponse{}, response.NewInternalServerErr(err.Error(), err)
	}

	return response.BoolResponse{
		Data: true,
	}, nil
}

func (r *roomService) JoinRoom(ctx context.Context, roomId uint64, userId uint64) (response.BoolResponse, error) {
	roomModel := &model.UserRoom{
		UserID: userId,
		RoomID: roomId,
		Role:   "member",
	}

	if err := r.roomRepo.JoinRoom(ctx, roomModel); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BoolResponse{Data: true}, response.NewNotFoundErr("room not found", err)
		}

		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return response.BoolResponse{Data: true}, response.NewBadRequestErr("user already join", err)
		}

		return response.BoolResponse{}, response.NewInternalServerErr(err.Error(), err)
	}

	return response.BoolResponse{
		Data: true,
	}, nil
}

func (r *roomService) DeleteRoom(ctx context.Context, roomId uint64) (response.BoolResponse, error) {
	if err := r.roomRepo.DeleteRoom(ctx, uint64(roomId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BoolResponse{Data: true}, response.NewNotFoundErr("room not found", err)
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

		var replyMsg *response.GetMessageByIdResponse
		if message.ReplyID != nil && message.ReplyContent != nil {
			replyMsg = &response.GetMessageByIdResponse{
				ID:      *message.ReplyID,
				Content: json.RawMessage(*message.ReplyContent),
			}
		}

		res = append(res, response.RoomMessageResponse{
			ID:           message.ID,
			Sender:       sender,
			ReplyContent: replyMsg,
			Content:      json.RawMessage(message.Content),
			CreatedAt:    message.CreatedAt,
			UpdatedAt:    message.UpdatedAt,
		})
	}

	return res, nil
}

func (r *roomService) CreateRoom(ctx context.Context, userId uint64, req *request.CreateRoomRequest) (response.GetRoomResponse, error) {
	if req == nil {
		return response.GetRoomResponse{}, response.NewBadRequestErr("invalid create room request", nil)
	}

	roomId := time.Now().UnixMicro()
	room := &model.Room{
		ID:        uint64(roomId),
		CreatorID: uint64(userId),
		Name:      req.Name,
		ImgUrl:    nil,
		CreatedAt: time.Now(),
	}
	if err := r.roomRepo.CreateRoom(ctx, room); err != nil {
		return response.GetRoomResponse{}, err
	}

	return response.GetRoomResponse{
		ID:        uint64(roomId),
		Name:      room.Name,
		ImgUrl:    room.ImgUrl,
		CreatedAt: room.CreatedAt}, nil
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
			Name:   room.Creator.Username,
			ImgUrl: room.Creator.ImgUrl,
		},
		CreatedAt: room.CreatedAt,
	}, nil
}
