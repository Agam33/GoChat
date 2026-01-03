package websocket

type BroadcastMessage struct {
	RoomId int64
	Data   []byte
}

type RoomHub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMessage
}

func NewRoomHub() *RoomHub {
	return &RoomHub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan BroadcastMessage, 512),
	}
}

func (h *RoomHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)

			if len(h.clients) == 0 {
				return
			}

		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- msg.Data:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
