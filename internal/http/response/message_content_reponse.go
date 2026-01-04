package response

import "time"

type TextContentReponse struct {
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type ImageContentResponse struct {
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	Images    []string  `json:"images"`
	CreatedAt time.Time `json:"createdAt"`
}
