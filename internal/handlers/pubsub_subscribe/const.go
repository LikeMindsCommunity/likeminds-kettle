package pubsub_subscribe

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
	TopicTypeChatroom = "chatroom"
	ParamTopic        = "topic"
)

const (
	ErrorUnableToCloseWs         = "unable to close ws error:"
	ErrorFailedUpgrader          = "failed to upgrade connection: %v"
	ErrorFailedDial              = "failed to dial connection: %v"
	ErrorPingSendClient          = "error sending ping to client: %v"
	ErrorPingSendWs              = "error sending ping to websocket server: %v"
	ErrorReadClientWs            = "error reading message from client: %v"
	ErrorWriteClientWs           = "error writing message to client: %v"
	ErrorReadServerWs            = "error reading message from server: %v"
	ErrorWriteServerWs           = "error writing message to server: %v"
	ErrorTopicMissing            = "topic is missing from request"
	ErrorTopicInvalid            = "invalid format of topic"
	ErrorUserUUIDMissing         = "user UUID is missing in header"
	ErrorUserChatroomAccess      = "unable to subscribe to chatroom - %s"
	ErrorUnmarshalErrorJson      = "unmarshal error: %v"
	ErrorTopicIDMissing          = "topic ID is missing in request"
	ErrorChatroomResponseInvalid = "invalid channel_details_data key in response"
)
