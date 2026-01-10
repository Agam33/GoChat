package chat

import (
	"context"
	"encoding/json"
	"errors"
	"go-chat/internal/http/response"
	"go-chat/internal/model"
	"go-chat/internal/websocket/event"
	"time"

	"gorm.io/gorm"
)

type ChatService interface {
	ReplyMessage(ctx context.Context, req *event.SendReplyEvent, contentData event.TextContentData) error
	DeleteMessage(ctx context.Context, req *event.DeleteMessageEvent) (response.BoolResponse, error)
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

func (c *chatService) ReplyMessage(ctx context.Context, req *event.SendReplyEvent, contentData event.TextContentData) error {
	jsonb, _ := json.Marshal(contentData)

	msgModel := &model.Message{
		ID:       uint64(time.Now().Unix()),
		SenderID: req.SenderId,
		ReplyID:  &req.ReplyTo,
		Content:  jsonb,
	}

	if err := c.chatRepo.SaveMessage(ctx, msgModel); err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	return nil
}

func (c *chatService) SaveTextMessage(ctx context.Context, senderId uint64, req event.TextContentData) error {
	jsonb, _ := json.Marshal(req)

	msgModel := &model.Message{
		ID:          uint64(time.Now().Unix()),
		SenderID:    senderId,
		ContentType: req.ContentType,
		Content:     jsonb,
	}

	if err := c.chatRepo.SaveMessage(ctx, msgModel); err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	return nil
}

func (c *chatService) DeleteMessage(ctx context.Context, req *event.DeleteMessageEvent) (response.BoolResponse, error) {
	err := c.chatRepo.DeleteMessage(ctx, req.SenderId, req.MessageId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.BoolResponse{}, response.NewNotFoundErr("message not found", err)
		}
		return response.BoolResponse{}, response.NewInternalServerErr("error delete message", err)
	}

	return response.BoolResponse{
		Data: true,
	}, nil
}
