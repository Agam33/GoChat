package response

type MessageTextResponse struct {
	ID      uint64             `json:"id"`
	Sender  UserReponse        `json:"sender"`
	Content TextContentReponse `json:"content"`
}

type MessageImgResponse struct {
	ID      uint64               `json:"id"`
	Sender  UserReponse          `json:"sender"`
	Content ImageContentResponse `json:"content"`
}
