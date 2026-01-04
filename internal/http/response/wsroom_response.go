package response

import (
	"encoding/json"
)

type WsChatResponse struct {
	Type   string          `json:"type"`
	RoomId int64           `json:"roomId"`
	Data   json.RawMessage `json:"data"`
}
