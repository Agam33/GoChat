package chat

import (
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
)

type ChatService interface {
	SaveMessage(roomId uint64, req request.SendTextRequest) (*response.WsChatResponse, error)
}

type chatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (c *chatService) SaveMessage(roomId uint64, req request.SendTextRequest) (*response.WsChatResponse, error) {
	
	return nil
}
