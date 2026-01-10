package event

type SendReplyEvent struct {
	RoomId   uint64 `json:"roomId"`
	SenderId uint64 `json:"senderId"`
	ReplyTo  uint64 `json:"replyTo"`
	Text     string `json:"text"`
}
