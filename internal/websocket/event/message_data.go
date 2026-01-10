package event

import "encoding/json"

type MessageData struct {
	RoomId  uint64          `json:"roomId"`
	Sender  ClienData       `json:"sender"`
	Content json.RawMessage `json:"content"`
}
