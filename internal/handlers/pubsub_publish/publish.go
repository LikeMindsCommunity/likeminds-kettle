package pubsub_publish

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// PublishConversationOnTopicTypeChatroom to publish Conversation on TopicTypeChatroomDynamic
func PublishConversationOnTopicTypeChatroom(c *gin.Context, chatroomID interface{}, userId string, deviceID string, response map[string]interface{}) {
	headers := utils.CreateHeadersFromToken(c, userId, deviceID)
	topicChatroom := fmt.Sprintf(TopicTypeChatroomDynamic, chatroomID)
	params := map[string]string{
		ParamTopicMessageType: TopicMessageTypeConversation,
	}

	publishDataOnPandemonium(topicChatroom, headers, params, response)
}

// PublishConversationOnTopicTypeCommunity to publish Conversation on TopicTypeCommunityDynamic
func PublishConversationOnTopicTypeCommunity(c *gin.Context, userId string, deviceID string, response map[string]interface{}) {
	headers := utils.CreateHeadersFromToken(c, userId, deviceID)
	var communityID = response["conversation"].(map[string]interface{})["member"].(map[string]interface{})["sdk_client_info"].(map[string]interface{})["community"]
	topicCommunity := fmt.Sprintf(TopicTypeCommunityDynamic, communityID)
	publishParams := map[string]string{
		ParamTopicMessageType: TopicMessageTypeConversation,
	}

	publishDataOnPandemonium(topicCommunity, headers, publishParams, response)
}

// PublishDROnTopicTypeChatroom handles calling the Pandemonium publish API for inactive user flow in delivered report
func PublishDROnTopicTypeChatroom(c *gin.Context, userId string) {
	deviceID := user.GetRequestingUserDeviceId(c)
	headers := utils.CreateHeadersFromToken(c, userId, deviceID)
	topicChatroom := fmt.Sprintf(TopicTypeChatroomDynamic, c.Query(ParamChatroomId))
	params := map[string]string{
		ParamTopicMessageType: TopicMessageDR,
	}

	// Prepare the response (body) for the inactive user case
	response := map[string]interface{}{
		"delivered_report": map[string]interface{}{
			"min_timestamp": c.Query(ParamMinTimeStamp),
			"max_timestamp": c.Query(ParamMaxTimeStamp),
		},
	}

	// Call the Pandemonium Publish API using your pre-defined function
	publishDataOnPandemonium(topicChatroom, headers, params, response)
}

func publishDataOnPandemonium(topicChatroom string, headers map[string]interface{}, params map[string]string, response map[string]interface{}) {
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.PandemoniumService, fmt.Sprintf(PublishEndPoint, topicChatroom), utils.POSTRequestRawBody, headers, params, response)
	if err != nil || statusCode != 200 {
		logging.Error(fmt.Sprintf("Error in publishing data on pandemonium: %v | response %v", err, string(respBytes)))
	}
}
