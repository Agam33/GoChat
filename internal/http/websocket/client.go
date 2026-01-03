package websocket

type Client struct {
	ClientId int64
	Send     chan []byte
}

func readPump() {}

func writePump() {}

func ServeWS()
