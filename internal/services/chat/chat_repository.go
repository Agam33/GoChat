package chat

import (
	"context"
	"go-chat/internal/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
	DeleteMessage(context.Context, uint64, uint64) error
	SaveMessage(context.Context, *model.Message) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (repo *chatRepository) DeleteMessage(ctx context.Context, userId uint64, msgId uint64) error {
	err := repo.db.WithContext(ctx).Delete(&model.Message{ID: msgId, SenderID: userId}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepository) SaveMessage(ctx context.Context, message *model.Message) error {
	return repo.db.WithContext(ctx).Create(message).Error
}
