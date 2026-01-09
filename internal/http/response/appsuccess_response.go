package response

import "go-chat/internal/utils/types"

type SuccessReponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type SuccessReponseWithMeta[T any] struct {
	Message string     `json:"message"`
	Data    T          `json:"data"`
	Meta    types.Meta `json:"meta"`
}
