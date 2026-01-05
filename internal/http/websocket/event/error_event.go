package event

import (
	"encoding/json"
	"errors"
	"go-chat/internal/http/response"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsErrorEvent struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendWsError(conn *websocket.Conn, err error) {
	code := http.StatusInternalServerError
	msg := "internal server error"

	var apperr *response.AppErr
	if errors.As(err, &apperr) {
		code = apperr.Code

		if code < 500 && code >= 400 {
			msg = apperr.Message
		}

		if code >= 500 && apperr.Err != nil {
			log.Printf("error wesocket: %v", err)
		}
	}

	payload, _ := json.Marshal(WsErrorEvent{
		Type:    "error",
		Code:    code,
		Message: msg,
	})

	frame, _ := json.Marshal(WSMessageEvent{
		Type: "error",
		Data: payload,
	})

	_ = conn.WriteMessage(websocket.TextMessage, frame)
}

func FataError(err error) bool {
	var apperr *response.AppErr
	if errors.As(err, &apperr) {
		return apperr.Code == http.StatusUnauthorized
	}
	return false
}
