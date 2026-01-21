package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"go-chat/internal/constant"
	"go-chat/internal/services/chat"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type ChatConsumerHandler struct {
	chatService chat.ChatService
}

func NewChatConsumerHandler(chatService chat.ChatService) *ChatConsumerHandler {
	return &ChatConsumerHandler{
		chatService: chatService,
	}
}

func (c *ChatConsumerHandler) Dispatch(ctx context.Context, msg amqp.Delivery) {
	var err error
	var retry bool

	switch msg.RoutingKey {
	case constant.MQRoutingChatSave:
		retry, err = c.handleChatSave(ctx, msg)
	case constant.MQRoutingChatReply:
		retry, err = c.handleChatReply(ctx, msg)
	default:
		msg.Nack(false, false)
		return
	}

	if err != nil {
		log.Printf("chat consumer err rk=%s retry=%v id=%s err=%v",
			msg.RoutingKey, retry, msg.MessageId, err)
		msg.Nack(false, retry)
		return
	}

	msg.Ack(false)
}

func (c *ChatConsumerHandler) handleChatSave(ctx context.Context, msg amqp.Delivery) (bool, error) {
	var evt SaveTextEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		return false, err
	}

	if err := c.chatService.SaveTextMessage(ctx, evt.UserID, evt.RoomID, evt.Content); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return false, err
		}
		return true, err
	}

	return false, nil
}

func (c *ChatConsumerHandler) handleChatReply(ctx context.Context, msg amqp.Delivery) (bool, error) {
	var evt SendReplyTextEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		return false, err
	}

	if err := c.chatService.ReplyMessage(ctx, &evt.SendReply, evt.Content, evt.ReplyMsg); err != nil {
		return true, err
	}

	return false, nil
}
