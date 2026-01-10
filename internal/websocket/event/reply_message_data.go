package event

import "encoding/json"

type ReplyTextData struct {
	ID           uint64          `json:"id"`
	RoomId       uint64          `json:"roomId"`
	ReplyContent json.RawMessage `json:"replyContent"`
	Content      json.RawMessage `json:"content"`
}
