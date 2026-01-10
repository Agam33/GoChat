package event

type DeleteMessageData struct {
	RoomId    uint64 `json:"roomdId"`
	MessageId uint   `json:"messageId"`
}
