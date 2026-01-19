package consumer

import (
	"go-chat/internal/http/response"
	"go-chat/internal/websocket/event"
)

type SaveTextEvent struct {
	UserID  uint64
	RoomID  uint64
	Content event.TextContentData
}

type SendReplyTextEvent struct {
	SendReply event.SendReplyEvent
	Content   event.TextContentData
	ReplyMsg  response.GetMessageByIdResponse
}
