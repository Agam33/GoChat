package event

import "time"

type TextContentData struct {
	ID          uint64    `json:"id"`
	ContentType string    `json:"contentType"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ImageContentData struct {
	ID          uint64    `json:"id"`
	ContentType string    `json:"contentType"`
	Text        string    `json:"text"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"createdAt"`
}
