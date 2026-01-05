package event

type SendTextEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
