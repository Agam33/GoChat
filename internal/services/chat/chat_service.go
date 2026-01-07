package chat

import (
	"context"
	"encoding/json"
	"go-chat/internal/model"
	"go-chat/internal/websocket/event"
)

type ChatService interface {
	SaveTextMessage(ctx context.Context, senderId uint64, req event.TextContentData) error
}

type chatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (c *chatService) SaveTextMessage(ctx context.Context, senderId uint64, req event.TextContentData) error {
	jsonb, _ := json.Marshal(req)

	msgModel := &model.Message{
		ID:       req.ID,
		SenderID: senderId,
		Type:     req.ContentType,
		Content:  jsonb,
	}

	if err := c.chatRepo.SaveMessage(ctx, msgModel); err != nil {
		return nil
	}

	return nil
}
