package pubsubCommon

import "time"

const (
	ReadBufferSizeDefault  = 4096
	WriteBufferSizeDefault = 4096
	PongWait               = 60 * time.Second
	// PingPeriod should be less than PongWait
	PingPeriod                = (PongWait * 9) / 10
	WriteWait                 = 10 * time.Second
	WsConnectionEstablished   = "connected to websocket server"
	PongReceivedPandemoniumWs = "received pong from pandemonium websocket server"
	PongReceivedClient        = "received pong from client"
	PingReceivedClient        = "received ping from client"
	PingReceivedWs            = "received ping from websocket server"
	PingSendClient            = "sending ping to client"
	PingSendWs                = "sending ping to websocket server"
	ReceivedMessageClientWs   = "received message from client having message type as %v"
	ReceivedMessageServerWs   = "received message from websocket server having message type as %v"
	ConnectionClosed          = "connection closed"
)

const (
	ParamTopic                                = "topic"
	TopicTypeChatroom                         = "chatroom"
	TopicTypeCommunity                        = "community"
	TopicMessageTypeCreateConversationRequest = "message.create.request"
	TopicTypeChatroomDynamic                  = "chatroom:%v"
	TopicTypeCommunityDynamic                 = "community:%v"
	TopicMessageTypeConversation              = "conversation"
	TopicMessageDR                            = "delivered_dr"
	ParamTopicMessageType                     = "topic_message_type"
	ParamChatroomId                           = "chatroom_id"
	ParamPage                                 = "page"
	ParamPageSize                             = "page_size"
	ChatroomAPIVersion                        = "1"
	ChatroomParticipantsAPIPage               = "0"
	ChatroomParticipantsAPIPageSize           = "100"
	ChatroomParticipantsAPIVersion            = "1"
	ChatroomParticipantsAPIPlatformCode       = "an"
	ChatroomParticipantsAPIVersionCode        = "210"
	ParamTotalParticipantsCountType           = "total_participants_count"
	ParamParticipantsType                     = "participants"
	ParamIsSecret                             = "is_secret"
	PublishEndPoint                           = "/publish/%s"
	ParamMinTimeStamp                         = "min_timestamp"
	ParamMaxTimeStamp                         = "max_timestamp"
	ParamCommunityID                          = "community_id"
)

const (
	ErrorUnableToCloseWs         = "unable to close ws error: %v"
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
	ErrorUserChatroomAccess      = "unable to access chatroom - %s"
	ErrorUnmarshalErrorJson      = "unmarshal error: %v"
	ErrorTopicIDMissing          = "topic ID is missing in request"
	ErrorChatroomResponseInvalid = "invalid channel_details_data key in response"
	ErrorMarshalErrorJson        = "marshal error: %v"
	ErrorWriteDeadlineWs         = "error while setting WriteDeadline on websocket:"
	ErrorReadDeadlineWs          = "error while setting ReadDeadline on websocket:"
)
