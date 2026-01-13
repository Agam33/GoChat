package response

import (
	"encoding/json"
)

type SaveTextMessageResponse struct {
	ID uint64 `json:"id"`
}

type GetMessageByIdResponse struct {
	ID      uint64          `json:"id"`
	Content json.RawMessage `json:"content"`
}
