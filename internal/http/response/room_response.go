package response

import (
	"encoding/json"
	"time"
)

type GetRoomResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	ImgUrl    *string   `json:"imgUrl"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetDetailRoomResponse struct {
	ID        uint64       `json:"id"`
	Name      string       `json:"name"`
	ImgUrl    *string      `json:"imgUrl"`
	Creator   UserResponse `json:"creator"`
	CreatedAt time.Time    `json:"createdAt"`
}

type RoomMessageResponse struct {
	ID          uint64          `json:"id"`
	ContentType string          `json:"contentType"`
	Sender      UserResponse    `json:"sender"`
	Content     json.RawMessage `json:"content"`
	CreatedAt   time.Time       `json:"createdAt"`
}
