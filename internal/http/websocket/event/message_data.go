package event

type MessageTextData struct {
	ID      uint64          `json:"id"`
	Sender  ClienData       `json:"sender"`
	Content TextContentData `json:"content"`
}

type MessageImgData struct {
	ID      uint64           `json:"id"`
	Sender  ClienData        `json:"sender"`
	Content ImageContentData `json:"content"`
}
