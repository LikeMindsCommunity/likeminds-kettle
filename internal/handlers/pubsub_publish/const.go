package pubsub_publish

const (
	TopicTypeChatroomDynamic     = "chatroom:%s"
	TopicTypeCommunityDynamic    = "community:%v"
	TopicMessageTypeConversation = "conversation"
	TopicMessageDR               = "delivered_dr"
)

const (
	PublishEndPoint       = "/publish/%s"
	ParamTopicMessageType = "topic_message_type"
	ParamChatroomId       = "chatroom_id"
	ParamMinTimeStamp     = "min_timestamp"
	ParamMaxTimeStamp     = "max_timestamp"
)
