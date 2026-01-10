package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/services/chat"
	"go-chat/internal/services/room"
	"go-chat/internal/services/user"
	"go-chat/internal/websocket/event"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler interface {
	ServeWS(*gin.Context)
}

type wsHandler struct {
	hub         *Hub
	upgrader    websocket.Upgrader
	userService user.UserService
	roomService room.RoomService
	chatService chat.ChatService
}

func NewWSHandler(hub *Hub, userService user.UserService, roomService room.RoomService, chatService chat.ChatService) WsHandler {
	return &wsHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		userService: userService,
		roomService: roomService,
		chatService: chatService,
	}
}

func (h *wsHandler) Disconnect(c *Client) {
	log.Printf("unregister: %v", c)
	h.hub.unregister <- c
}

func (h *wsHandler) Dispatch(c *Client, msg event.WSMessageEvent) error {
	switch msg.Action {
	case "room.join":
		var roomJoin event.RoomEvent
		if err := json.Unmarshal(msg.Data, &roomJoin); err != nil {
			return response.NewBadRequestErr("can't parse request room join", err)
		}

		h.hub.subscribe <- &Subscriber{
			C:     c,
			Topic: BuildWSTopic("room", "chat", roomJoin.RoomId),
		}

		content := event.TextContentData{
			ContentType: "system",
			Text:        fmt.Sprintf("%s joined.", c.displayName),
		}

		if err := h.sendSystemChat(c, uint64(roomJoin.RoomId), content); err != nil {
			return err
		}

	case "room.leave": // temporary leave
		var roomLeave event.RoomEvent
		if err := json.Unmarshal(msg.Data, &roomLeave); err != nil {
			return response.NewBadRequestErr("can't parse request room join", err)
		}

		h.hub.unsubscribe <- &Subscriber{
			Topic: BuildWSTopic("room", "chat", roomLeave.RoomId),
			C:     c,
		}
	case "room.send.text":
		var sendText event.SendTextEvent
		if err := json.Unmarshal(msg.Data, &sendText); err != nil {
			return response.NewBadRequestErr("can't parse request chat text", err)
		}

		if err := h.roomSendText(c, sendText); err != nil {
			return err
		}
	case "room.reply.text":
		var replyMsg event.SendReplyEvent
		if err := json.Unmarshal(msg.Data, &replyMsg); err != nil {
			return response.NewBadRequestErr("can't parse request reply chat text", err)
		}

		if err := h.sendReplyText(c, &replyMsg); err != nil {
			return err
		}
	case "room.delete.message":
		var delMsg event.DeleteMessageEvent
		if err := json.Unmarshal(msg.Data, &delMsg); err != nil {
			return response.NewBadRequestErr("can't parse request delete message", err)
		}

		if err := h.deleteMessage(&delMsg); err != nil {
			return err
		}

	default:
	}

	return nil
}

func (h *wsHandler) sendSystemChat(c *Client, roomId uint64, content event.TextContentData) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := h.chatService.SaveTextMessage(ctx, c.UserId, content); err != nil {
		return err
	}

	contentData, _ := json.Marshal(content)

	evt := &event.WSMessageEvent{
		Action: "room.system.text",
		Data:   contentData,
	}
	evtData, _ := json.Marshal(evt)

	h.hub.broadcast <- BroadcastMessage{
		Topic: BuildWSTopic("room", "chat", int64(roomId)),
		Data:  evtData,
	}

	return nil
}

func (h *wsHandler) deleteMessage(delMsg *event.DeleteMessageEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := h.chatService.DeleteMessage(ctx, delMsg)
	if err != nil {
		return err
	}

	msgText := event.DeleteMessageData{
		RoomId:    delMsg.RoomId,
		MessageId: uint(delMsg.MessageId),
	}

	rawContent, err := json.Marshal(msgText)
	if err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	evt := event.WSMessageEvent{
		Action: "room.delete.message",
		Data:   rawContent,
	}
	evtData, _ := json.Marshal(evt)

	h.hub.broadcast <- BroadcastMessage{
		Topic: BuildWSTopic("room", "chat", int64(delMsg.RoomId)),
		Data:  evtData,
	}

	return nil
}

func (h *wsHandler) sendReplyText(c *Client, sendReply *event.SendReplyEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	contentData := event.TextContentData{
		ContentType: "text",
		Text:        sendReply.Text,
	}

	rawContent, err := json.Marshal(contentData)
	if err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	if err := h.chatService.ReplyMessage(ctx, sendReply, contentData); err != nil {
		return err
	}

	msgText := event.MessageData{
		RoomId: sendReply.RoomId,
		Sender: event.ClienData{
			ID:     c.UserId,
			Name:   c.displayName,
			ImgUrl: &c.avatarURL,
		},
		Content: rawContent,
	}
	textData, _ := json.Marshal(msgText)

	evt := event.WSMessageEvent{
		Action: "room.chat.reply",
		Data:   textData,
	}
	rawEvt, _ := json.Marshal(evt)

	h.hub.broadcast <- BroadcastMessage{
		Topic: BuildWSTopic("room", "chat", int64(sendReply.RoomId)),
		Data:  rawEvt,
	}

	return nil
}

func (h *wsHandler) roomSendText(c *Client, sendText event.SendTextEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	contentData := event.TextContentData{
		ContentType: "text",
		Text:        sendText.Text,
	}

	rawContent, err := json.Marshal(contentData)
	if err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	err = h.chatService.SaveTextMessage(ctx, c.UserId, contentData)
	if err != nil {
		return err
	}

	msgText := &event.MessageData{
		RoomId: uint64(sendText.RoomId),
		Sender: event.ClienData{
			ID:     c.UserId,
			Name:   c.displayName,
			ImgUrl: &c.avatarURL,
		},
		Content: rawContent,
	}
	textData, _ := json.Marshal(msgText)

	evt := &event.WSMessageEvent{
		Action: "room.chat.send",
		Data:   textData,
	}
	evtData, _ := json.Marshal(evt)

	h.hub.broadcast <- BroadcastMessage{
		Topic: BuildWSTopic("room", "chat", sendText.RoomId),
		Data:  evtData,
	}

	return nil
}

func (h *wsHandler) ServeWS(c *gin.Context) {
	userId := c.GetUint64(constant.CtxUserIDKey)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("can't serve websocket: %v", err)
		c.Error(response.NewInternalServerErr("error upgrader in ServeChatWs", err))
		return
	}

	usr, err := h.userService.GetById(c.Request.Context(), userId)
	if err != nil {
		c.Error(err)
		return
	}

	client := &Client{
		UserId:      usr.ID,
		Conn:        conn,
		Send:        make(chan []byte, 512),
		displayName: usr.Name,
		avatarURL:   *usr.ImgUrl,
	}

	h.hub.register <- client

	go client.ReadPump(c.Request.Context(), h)
	go client.WritePump()
}
