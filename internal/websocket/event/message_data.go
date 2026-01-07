package event

type MessageTextData struct {
	RoomId  uint64          `json:"roomId"`
	Sender  ClienData       `json:"sender"`
	Content TextContentData `json:"content"`
}

type MessageImgData struct {
	RoomId  uint64           `json:"roomId"`
	Sender  ClienData        `json:"sender"`
	Content ImageContentData `json:"content"`
}
