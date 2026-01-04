package websocket

import (
	"context"
	"encoding/json"
	"go-chat/internal/http/request"
	"go-chat/internal/http/response"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 15 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	RoomId int64
	UserId uint64
	Name   string
	ImgUrl *string
	Send   chan []byte
	Conn   *websocket.Conn
	Ctx    context.Context
}

func (c *Client) ReadPump(roomHub *RoomHub) {
	defer func() {
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msgResp response.WsChatResponse
		if err := c.Conn.ReadJSON(&msgResp); err != nil {
			log.Printf("error msgReponse %v", err)
			break
		}

		switch msgResp.Type {
		case "chat.send.text":
			var sendTextRequest request.SendTextRequest
			if err := json.Unmarshal(msgResp.Data, &sendTextRequest); err != nil {
				log.Printf("error send request %v", err)
				continue
			}

		case "chat.reply":
		case "chat.join":
		case "room.leave":

		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				)
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
