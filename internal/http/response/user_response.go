package response

type UserResponse struct {
	ID     uint64  `json:"id"`
	Name   string  `json:"name"`
	ImgUrl *string `json:"imgUrl"`
}
