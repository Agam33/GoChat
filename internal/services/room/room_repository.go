package room

import (
	"context"
	"go-chat/internal/model"

	"gorm.io/gorm"
)

type RoomRepository interface {
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
