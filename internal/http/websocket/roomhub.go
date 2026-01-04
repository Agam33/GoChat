package websocket

import (
	"go-chat/internal/http/request"
	"go-chat/internal/model"
	"go-chat/internal/services/chat"
	"go-chat/internal/services/room"
	"log"
)

type BroadcastMessage struct {
	RoomId int64
	Data   []byte
}

type RoomHub struct {
	Rooms      map[int64]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage

	chatService chat.ChatService
	roomService room.RoomService
}

func NewRoomHub(chatService chat.ChatService, roomService room.RoomService) *RoomHub {
	return &RoomHub{
		Rooms:      make(map[int64]map[*Client]bool),
		Register:   make(chan *Client, 256),
		Unregister: make(chan *Client, 256),
		Broadcast:  make(chan BroadcastMessage, 512),

		chatService: chatService,
		roomService: roomService,
	}
}

func (h *RoomHub) SendMessageText(roomId int64, req request.SendTextRequest) {
	err := h.chatService.SaveMessage(uint64(roomId), req)
	if err != nil {
		log.Printf("error in send message %v", err)
		return
	}
}

func (h *RoomHub) GetRoom(roomId int64) (*model.Room, error) {
	return 0, nil
}

func (h *RoomHub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.RoomId] == nil {
				h.Rooms[client.RoomId] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomId][client] = true

			go client.ReadPump(h)
			go client.WritePump()

		case client := <-h.Unregister:
			clients := h.Rooms[client.RoomId]
			delete(clients, client)

		case msg := <-h.Broadcast:
			clients := h.Rooms[msg.RoomId]
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
