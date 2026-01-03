package response

import "time"

type ChatResponse struct {
	Type      string    `json:"type"`
	RoomId    int64     `json:"roomId"`
	Data      any       `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
}
