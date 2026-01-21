package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-chat/internal/constant"
	"go-chat/internal/http/response"
	"go-chat/internal/rabbitmq"
	"go-chat/internal/services/chat"
	chatCmsr "go-chat/internal/services/chat/consumer"
	"go-chat/internal/services/room"
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
	publisher   rabbitmq.Publisher
	roomService room.RoomService
	chatService chat.ChatService
}

func NewWSHandler(hub *Hub, publisher rabbitmq.Publisher, roomService room.RoomService, chatService chat.ChatService) WsHandler {
	return &wsHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		publisher:   publisher,
		roomService: roomService,
		chatService: chatService,
	}
}

func (h *wsHandler) Disconnect(c *Client) {
	log.Printf("unregister: %v", c)
	h.hub.unregister <- c
}

func (h *wsHandler) Dispatch(c *Client, msg event.WSMessageEvent) bool {
	switch msg.Action {
	case "room_join":
		var roomJoin event.RoomEvent
		if err := json.Unmarshal(msg.Data, &roomJoin); err != nil {
			h.sendWsError(c, response.NewBadRequestErr("can't parse request room join", err))
			return false
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		h.hub.subscribe <- &Subscriber{
			C:     c,
			Topic: BuildWSTopic("room", "chat", roomJoin.RoomId),
		}

		_, err := h.roomService.JoinRoom(ctx, uint64(roomJoin.RoomId), c.UserId)
		if err != nil {
			return true
		}

		content := event.TextContentData{
			ContentType: "system",
			Text:        fmt.Sprintf("%s joined.", c.displayName),
			CreatedAt:   time.Now(),
		}

		if err := h.sendSystemChat(ctx, c, uint64(roomJoin.RoomId), content); err != nil {
			h.sendWsError(c, err)
			return false
		}

	case "room_leave":
		var roomLeave event.RoomEvent
		if err := json.Unmarshal(msg.Data, &roomLeave); err != nil {
			h.sendWsError(c, response.NewBadRequestErr("can't parse request room join", err))
			return false
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := h.roomService.LeaveRoom(ctx, uint64(roomLeave.RoomId), c.UserId); err != nil {
			h.sendWsError(c, err)
			return false
		}

		content := event.TextContentData{
			ContentType: "system",
			Text:        fmt.Sprintf("%s left.", c.displayName),
			CreatedAt:   time.Now(),
		}

		if err := h.sendSystemChat(ctx, c, uint64(roomLeave.RoomId), content); err != nil {
			h.sendWsError(c, err)
			return false
		}

		h.hub.unsubscribe <- &Subscriber{
			Topic: BuildWSTopic("room", "chat", roomLeave.RoomId),
			C:     c,
		}
	case "room_send_text":
		var sendText event.SendTextEvent
		if err := json.Unmarshal(msg.Data, &sendText); err != nil {
			h.sendWsError(c, response.NewBadRequestErr("can't parse request chat text", err))
			return false
		}

		if err := h.roomSendText(c, sendText); err != nil {
			h.sendWsError(c, err)
			return false
		}
	case "room_reply_text":
		var replyMsg event.SendReplyEvent
		if err := json.Unmarshal(msg.Data, &replyMsg); err != nil {
			h.sendWsError(c, response.NewBadRequestErr("can't parse request reply chat text", err))
			return false
		}

		if err := h.sendReplyText(c, &replyMsg); err != nil {
			h.sendWsError(c, err)
			return false
		}
	case "room_delete_message":
		var delMsg event.DeleteMessageEvent
		if err := json.Unmarshal(msg.Data, &delMsg); err != nil {
			h.sendWsError(c, response.NewBadRequestErr("can't parse request delete message", err))
			return false
		}

		if err := h.deleteMessage(&delMsg); err != nil {
			h.sendWsError(c, err)
			return false
		}

	default:
	}

	return true
}

func (h *wsHandler) sendSystemChat(ctx context.Context, c *Client, roomId uint64, content event.TextContentData) error {
	dt, _ := json.Marshal(chatCmsr.SaveTextEvent{
		UserID:  c.UserId,
		RoomID:  roomId,
		Content: content,
	})
	h.publisher.Publish(ctx, constant.MQExchangeChat, constant.MQKindTopic, constant.MQRoutingChatSave, dt)

	contentData, _ := json.Marshal(content)

	evt := &event.WSMessageEvent{
		Action: "room_chat_system",
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

	if _, err := h.chatService.DeleteMessage(ctx, delMsg); err != nil {
		return response.NewNotFoundErr("message not found", err)
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
		Action: "room_delete_message",
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	replyMsg, err := h.chatService.GetMessageById(ctx, sendReply.ReplyTo)
	if err != nil {
		return err
	}

	contentData := event.TextContentData{
		Text:        sendReply.Text,
		ContentType: "text",
		CreatedAt:   time.Now(),
	}

	dt, _ := json.Marshal(chatCmsr.SendReplyTextEvent{
		SendReply: *sendReply,
		Content:   contentData,
		ReplyMsg:  replyMsg,
	})
	h.publisher.Publish(ctx, constant.MQExchangeChat, constant.MQKindTopic, constant.MQRoutingChatReply, dt)

	rawContent, err := json.Marshal(contentData)
	if err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	msgText := event.MessageData{
		RoomId: sendReply.RoomId,
		Sender: event.ClienData{
			ID:     c.UserId,
			Name:   c.displayName,
			ImgUrl: c.avatarURL,
		},
		ReplyContent: &event.ReplyContentData{
			ID:      replyMsg.ID,
			Content: replyMsg.Content,
		},
		Content:   rawContent,
		CreatedAt: contentData.CreatedAt,
	}
	textData, _ := json.Marshal(msgText)

	evt := event.WSMessageEvent{
		Action: "room_chat_reply",
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
		CreatedAt:   time.Now(),
	}

	rawContent, err := json.Marshal(contentData)
	if err != nil {
		return response.NewInternalServerErr(err.Error(), err)
	}

	dt, _ := json.Marshal(chatCmsr.SaveTextEvent{
		UserID:  c.UserId,
		RoomID:  uint64(sendText.RoomId),
		Content: contentData,
	})
	h.publisher.Publish(ctx, constant.MQExchangeChat, constant.MQKindTopic, constant.MQRoutingChatSave, dt)

	msgText := &event.MessageData{
		RoomId: uint64(sendText.RoomId),
		Sender: event.ClienData{
			ID:     c.UserId,
			Name:   c.displayName,
			ImgUrl: c.avatarURL,
		},
		Content:   rawContent,
		CreatedAt: contentData.CreatedAt,
		UpdatedAt: contentData.CreatedAt,
	}
	textData, _ := json.Marshal(msgText)

	evt := &event.WSMessageEvent{
		Action: "room_chat_send",
		Data:   textData,
	}
	evtData, _ := json.Marshal(evt)

	h.hub.broadcast <- BroadcastMessage{
		Topic: BuildWSTopic("room", "chat", sendText.RoomId),
		Data:  evtData,
	}

	return nil
}

func (h *wsHandler) sendWsError(c *Client, err error) {
	code := http.StatusInternalServerError
	msg := "internal server error"

	var apperr *response.AppErr
	if errors.As(err, &apperr) {
		code = apperr.Code

		if code < 500 && code >= 400 {
			msg = apperr.Message
		}

		if code >= 500 && apperr.Err != nil {
			log.Printf("error wesocket: %v", err)
		}
	} else {
		log.Printf("[INTERNAL SERVER ERROR] (ReadPump) %v", err)
		return
	}

	payload, _ := json.Marshal(event.WsErrorEvent{
		Type:    "error",
		Code:    code,
		Message: msg,
	})

	frame, _ := json.Marshal(event.WSMessageEvent{
		Action: "error",
		Data:   payload,
	})

	select {
	case c.Send <- frame:
	default:
		return
	}
}

func (h *wsHandler) ServeWS(c *gin.Context) {
	// prevent wrong client
	if !websocket.IsWebSocketUpgrade(c.Request) {
		c.AbortWithError(http.StatusBadRequest, response.NewBadRequestErr("can't upgrade to websocket", nil))
		return
	}

	ctxUser, exists := c.Get(constant.CtxUser)
	if !exists {
		c.AbortWithError(http.StatusUnauthorized, response.NewUnauthorized())
		return
	}

	usr, ok := ctxUser.(response.UserResponse)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, response.NewInternalServerErr("serve-ws: invalid user type", nil))
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, response.NewBadRequestErr("can't serve websocket", err))
		return
	}

	client := &Client{
		UserId:      usr.ID,
		Conn:        conn,
		Send:        make(chan []byte, 512),
		displayName: usr.Name,
		avatarURL:   usr.ImgUrl,
	}

	h.hub.register <- client

	go client.ReadPump(h)
	go client.WritePump()
}
