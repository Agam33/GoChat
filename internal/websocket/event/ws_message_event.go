package event

import (
	"encoding/json"
)

type WSMessageEvent struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}
