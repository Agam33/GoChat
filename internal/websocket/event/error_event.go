package event

type WsErrorEvent struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
