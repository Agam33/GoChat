package event

import (
	"encoding/json"
)

type WSMessageEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
