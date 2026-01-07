package handler

type RoomHandler interface {
}

type roomHandler struct{}

func NewRoomHandler() RoomHandler {
	return &roomHandler{}
}
