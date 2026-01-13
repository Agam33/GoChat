package event

import "encoding/json"

type ReplyContentData struct {
	ID      uint64          `json:"id"`
	Content json.RawMessage `json:"content"`
}
