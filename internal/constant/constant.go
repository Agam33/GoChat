package constant

const (
	StatusFailed  = "failed"
	StatusSuccess = "success"
	StatusError   = "error"

	CtxUserIDKey    = "userID"
	CtxUser         = "user"
	CtxRefreshToken = "refreshToken"

	Authorization = "Authorization"

	MQKindTopic = "topic"

	MQBindKeyChat = "chat.#"

	MQExchangeChat = "chat"

	QNameChat = "chat-service-q"

	MQRoutingChatSave  = "chat.save.text"
	MQRoutingChatDel   = "chat.del"
	MQRoutingChatReply = "chat.reply.text"
)
