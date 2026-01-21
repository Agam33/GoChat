package response

import (
	"encoding/json"
	"go-chat/internal/utils/types"
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
	ID           uint64                  `json:"id"`
	Sender       UserResponse            `json:"sender"`
	Content      json.RawMessage         `json:"content" swaggertype:"object"`
	ReplyContent *GetMessageByIdResponse `json:"replyContent"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

// for swagger doc
type GetRoomsResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Meta    types.Meta        `json:"meta"`
	Data    []GetRoomResponse `json:"data"`
}
