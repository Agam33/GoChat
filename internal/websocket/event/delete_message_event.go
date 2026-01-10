package event

type DeleteMessageEvent struct {
	RoomId    uint64 `json:"roomId"`
	SenderId  uint64 `json:"senderId"`
	MessageId uint64 `json:"messageId"`
}
