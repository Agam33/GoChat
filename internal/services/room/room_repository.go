package room

import (
	"context"
	"errors"
	"go-chat/internal/model"
	"go-chat/internal/utils"
	"go-chat/internal/utils/types"

	"gorm.io/gorm"
)

type RoomRepository interface {
	DeleteRoom(context.Context, uint64) error
	GetRoomMessages(context.Context, uint64, *types.Pagination) ([]model.Message, error)
	CreateRoom(context.Context, *model.Room) error
	GetById(context.Context, uint64) (*model.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func (r *roomRepository) DeleteRoom(ctx context.Context, roomId uint64) error {
	if err := r.db.WithContext(ctx).Delete(&model.Room{ID: roomId}).Error; err != nil {
		return err
	}

	return nil
}

func (r *roomRepository) GetRoomMessages(ctx context.Context, roomId uint64, pagination *types.Pagination) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.WithContext(ctx).Where("room_id = ?", roomId).
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(utils.PageOffset(pagination.Page, pagination.Limit)).
		Preload("Sender").Find(&messages).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []model.Message{}, nil
		}
		return []model.Message{}, err
	}

	return messages, nil
}

func (r *roomRepository) CreateRoom(ctx context.Context, room *model.Room) error {
	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) GetById(ctx context.Context, roomId uint64) (*model.Room, error) {
	var room model.Room
	if err := r.db.WithContext(ctx).
		Preload("Creator").
		First(&room, roomId).Error; err != nil {
		return nil, err
	}

	return &room, nil
}
