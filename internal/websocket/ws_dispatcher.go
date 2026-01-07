package websocket

import (
	"go-chat/internal/websocket/event"
)

type Dispatcher interface {
	Dispatch(c *Client, msg event.WSMessageEvent) error
	Disconnect(c *Client)
}
