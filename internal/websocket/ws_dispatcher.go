package websocket

import (
	"go-chat/internal/websocket/event"
)

type Dispatcher interface {
	Dispatch(c *Client, msg event.WSMessageEvent) bool
	Disconnect(c *Client)
}
