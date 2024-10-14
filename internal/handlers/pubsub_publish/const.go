package pubsub_publish

const (
	TopicTypeChatroomDynamic     = "chatroom:%v"
	TopicTypeCommunityDynamic    = "community:%v"
	TopicMessageTypeConversation = "conversation"
	TopicMessageDR               = "delivered_dr"
)

const (
	PublishEndPoint       = "/publish/%s"
	ParamTopicMessageType = "topic_message_type"
	ParamMinTimeStamp     = "min_timestamp"
	ParamMaxTimeStamp     = "max_timestamp"
)
