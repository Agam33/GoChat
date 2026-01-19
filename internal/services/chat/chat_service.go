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
	GetMessageById(ctx context.Context, msgId uint64) (response.GetMessageByIdResponse, error)
	ReplyMessage(ctx context.Context, req *event.SendReplyEvent, contentData event.TextContentData, replyMsg response.GetMessageByIdResponse) error
	DeleteMessage(ctx context.Context, req *event.DeleteMessageEvent) (response.BoolResponse, error)
	SaveTextMessage(ctx context.Context, senderId uint64, roomId uint64, req event.TextContentData) error
}

type chatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (c *chatService) GetMessageById(ctx context.Context, msgId uint64) (response.GetMessageByIdResponse, error) {
	msg, err := c.chatRepo.GetMessageById(ctx, msgId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.GetMessageByIdResponse{}, response.NewNotFoundErr("msg not found", err)
		}

		return response.GetMessageByIdResponse{}, response.NewInternalServerErr(err.Error(), err)
	}

	return response.GetMessageByIdResponse{
		ID:      msg.ID,
		Content: json.RawMessage(msg.Content),
	}, nil
}

func (c *chatService) ReplyMessage(ctx context.Context, req *event.SendReplyEvent, contentData event.TextContentData, resMessage response.GetMessageByIdResponse) error {
	if err := c.chatRepo.WithTransaction(ctx, func(chatRepo ChatRepository) error {
		replyId, err := chatRepo.SaveReplyMessage(ctx, contentData.ContentType, resMessage.Content)
		if err != nil {
			return response.NewInternalServerErr(err.Error(), err)
		}

		jsonb, _ := json.Marshal(contentData)
		msgModel := &model.Message{
			ID:          uint64(time.Now().UnixMilli()),
			RoomID:      req.RoomId,
			SenderID:    req.SenderId,
			ReplyID:     &replyId,
			ContentType: contentData.ContentType,
			Content:     jsonb,
			CreatedAt:   contentData.CreatedAt,
			UpdatedAt:   contentData.CreatedAt,
		}

		if err := chatRepo.SaveMessage(ctx, msgModel); err != nil {
			return response.NewInternalServerErr(err.Error(), err)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *chatService) SaveTextMessage(ctx context.Context, senderId uint64, roomId uint64, req event.TextContentData) error {
	jsonb, _ := json.Marshal(req)

	msgModel := &model.Message{
		ID:          uint64(time.Now().Unix()),
		RoomID:      roomId,
		SenderID:    senderId,
		ContentType: req.ContentType,
		Content:     jsonb,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.CreatedAt,
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
