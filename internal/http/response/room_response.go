package response

import "time"

type GetRoomResponse struct {
	ID        uint64       `json:"id"`
	Name      string       `json:"name"`
	ImgUrl    *string      `json:"imgUrl"`
	Creator   UserResponse `json:"creator"`
	CreatedAt time.Time    `json:"createdAt"`
}
