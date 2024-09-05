package pubsub

import "time"

const (
	ReadBufferSizeDefault   = 4096
	WriteBufferSizeDefault  = 4096
	PingPeriod              = ((60 * time.Second) * 9) / 10
	WsConnectionEstablished = "Connected to websocket server"
	PongReceivedWs          = "Received pong from websocket server"
	PongReceivedClient      = "Received pong from client"
	PingReceivedClient      = "Received ping from client"
	PingReceivedWs          = "Received ping from websocket server"
	PingSendClient          = "Sending ping to client"
	PingSendWs              = "Sending ping to websocket server"
	ReceivedMessageClientWs = "Received message from client having message type as %v"
	ReceivedMessageServerWs = "Received message from websocket server having message type as %v"
	ConnectionClosed        = "Connection closed"
)

const (
	TopicTypeChatroom            = "chatroom:%s"
	TopicTypeCommunity           = "community:%v"
	TopicMessageTypeConversation = "conversation"
)

const ParamTopicMessageType = "topic_message_type"

const (
	ErrorUnableToCloseWs = "Unable to close ws error:"
	ErrorFailedUpgrader  = "Failed to upgrade connection: %v"
	ErrorFailedDial      = "Failed to dial connection: %v"
	ErrorPingSendClient  = "Error sending ping to client: %v"
	ErrorPingSendWs      = "Error sending ping to websocket server: %v"
	ErrorReadClientWs    = "Error reading message from client: %v"
	ErrorWriteClientWs   = "Error writing message to client: %v"
	ErrorReadServerWs    = "Error reading message from server: %v"
	ErrorWriteServerWs   = "Error writing message to server: %v"
)
