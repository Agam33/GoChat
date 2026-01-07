package chat

import (
	"context"
	"go-chat/internal/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
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

func (repo *chatRepository) SaveMessage(ctx context.Context, message *model.Message) error {
	return repo.db.WithContext(ctx).Create(message).Error
}
