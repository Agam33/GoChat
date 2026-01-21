package request

type CreateRoomRequest struct {
	Name string `json:"name"`
}

type JoinRoomRequst struct {
	RoomId uint64 `json:"roomId"`
}
