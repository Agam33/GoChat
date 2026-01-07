package websocket

import (
	"context"
	"go-chat/internal/http/response"
	"go-chat/internal/websocket/event"
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
	UserId uint64
	Send   chan []byte
	Conn   *websocket.Conn

	displayName string
	avatarURL   string
}

func (c *Client) ReadPump(ctx context.Context, d Dispatcher) {
	defer func() {
		d.Disconnect(c)
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(appData string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msgResp event.WSMessageEvent
		if err := c.Conn.ReadJSON(&msgResp); err != nil {
			event.SendWsError(c.Conn, response.NewBadRequestErr("invalid message event", err))
			break
		}

		if err := d.Dispatch(c, msgResp); err != nil {
			event.SendWsError(c.Conn, err)
			continue
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
				log.Printf("error ws write: %v", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error write deadline %v", err)
				return
			}
		}
	}
}
