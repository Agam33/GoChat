package chat

import "go-chat/internal/http/request"

type ChatService interface {
	SaveMessage(roomId uint64, req request.SendTextRequest) error
}

type chatService struct {
	chatRepo ChatRepository
}

func NewChatService(chatRepo ChatRepository) ChatService {
	return &chatService{
		chatRepo: chatRepo,
	}
}

func (c *chatService) SaveMessage(roomId uint64, req request.SendTextRequest) error {
	return nil
}
