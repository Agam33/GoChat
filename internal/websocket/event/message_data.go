package event

import (
	"encoding/json"
	"time"
)

type MessageData struct {
	RoomId       uint64            `json:"roomId"`
	Sender       ClienData         `json:"sender"`
	ReplyContent *ReplyContentData `json:"replyContent"`
	Content      json.RawMessage   `json:"content"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updateAt"`
}
