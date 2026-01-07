package event

type SendTextEvent struct {
	RoomId int64  `json:"roomId"`
	Text   string `json:"text"`
}

type SendImageEvent struct {
	RoomId int64    `json:"roomId"`
	Text   string   `json:"text"`
	Images []string `json:"images"`
}
