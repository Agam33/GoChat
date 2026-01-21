package chat

import (
	"context"
	"go-chat/internal/model"

	"gorm.io/gorm"
)

type ChatRepository interface {
	WithTransaction(ctx context.Context, cb func(chatRepo ChatRepository) error) error
	GetMessageById(ctx context.Context, msgId uint64) (*model.Message, error)
	DeleteMessage(ctx context.Context, userId uint64, msgId uint64, content []byte) error
	SaveMessage(ctx context.Context, message *model.Message) error
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (repo *chatRepository) WithTransaction(ctx context.Context, cb func(chatRepo ChatRepository) error) error {
	if err := repo.db.Transaction(func(tx *gorm.DB) error {
		rp := &chatRepository{db: tx}
		if err := cb(rp); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (repo *chatRepository) GetMessageById(ctx context.Context, msgId uint64) (*model.Message, error) {
	var msg model.Message
	if err := repo.db.WithContext(ctx).Where("id = ?", msgId).First(&msg).Error; err != nil {
		return nil, err
	}

	return &msg, nil
}

func (repo *chatRepository) DeleteMessage(ctx context.Context, userId uint64, msgId uint64, content []byte) error {
	err := repo.db.WithContext(ctx).Model(&model.Message{}).Where("id = ? AND sender_id = ?", msgId, userId).Updates(map[string]any{
		"content_type": "text",
		"content":      content,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepository) SaveMessage(ctx context.Context, message *model.Message) error {
	return repo.db.WithContext(ctx).Create(message).Error
}
