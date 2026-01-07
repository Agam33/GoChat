package websocket

import "log"

type BroadcastMessage struct {
	Topic string
	Data  []byte
}

type Subscriber struct {
	C     *Client
	Topic string
}

type Hub struct {
	Topics map[string]map[*Client]struct{}

	register   chan *Client
	unregister chan *Client

	subscribe   chan *Subscriber
	unsubscribe chan *Subscriber

	broadcast chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		Topics:     make(map[string]map[*Client]struct{}),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
		broadcast:  make(chan BroadcastMessage, 512),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Printf("client regis: %v", client)

		case client := <-h.unregister:
			for topic, clients := range h.Topics {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.Topics, topic)
				}
			}

		case sub := <-h.subscribe:
			h.Topics[sub.Topic] = map[*Client]struct{}{}

		case sub := <-h.unsubscribe:
			clients := h.Topics[sub.Topic]
			delete(clients, sub.C)

		case msg := <-h.broadcast:
			clients := h.Topics[msg.Topic]
			for c := range clients {
				select {
				case c.Send <- msg.Data:
				default:
					close(c.Send)
					delete(clients, c)
				}
			}
		}
	}
}
