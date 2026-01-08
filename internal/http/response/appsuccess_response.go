package response

type SuccessReponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
