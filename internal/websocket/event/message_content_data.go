package event

import "time"

type TextContentData struct {
	ContentType string    `json:"contentType"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"-"`
}

type ImageContentData struct {
	ContentType string    `json:"contentType"`
	Text        string    `json:"text"`
	Images      []string  `json:"images"`
	CreatedAt   time.Time `json:"-"`
}
