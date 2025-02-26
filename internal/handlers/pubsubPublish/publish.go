package pubsubPublish

import (
	"fmt"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubCommon"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// PublishDROnTopicTypeChatroom handles calling the Pandemonium publish API for inactive user flow in delivered report
func PublishDROnTopicTypeChatroom(headers map[string]interface{}, minTimeStamp, maxTimeStamp string, chatroomID, communityID interface{}) {
	if chatroomID == nil {
		logging.Error("PublishDROnTopicTypeChatroom: chatroom ID is missing")
		return
	}
	if communityID == nil {
		logging.Error("PublishDROnTopicTypeChatroom: community ID is missing")
		return
	}

	topicChatroom := fmt.Sprintf(pubsubCommon.TopicTypeChatroomDynamic, chatroomID)
	params := map[string]string{
		pubsubCommon.ParamTopicMessageType: pubsubCommon.TopicMessageDR,
	}

	// Prepare the response (body) for the inactive user case
	response := map[string]interface{}{
		pubsubCommon.ParamMinTimeStamp: minTimeStamp,
		pubsubCommon.ParamMaxTimeStamp: maxTimeStamp,
		pubsubCommon.ParamCommunityID:  communityID,
	}

	// Call the Pandemonium Publish API using your pre-defined function
	publishDataOnPandemonium(topicChatroom, headers, params, response)
}

func publishDataOnPandemonium(topicChatroom string, headers map[string]interface{}, params map[string]string, response map[string]interface{}) {
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.PandemoniumService, fmt.Sprintf(pubsubCommon.PublishEndPoint, topicChatroom), utils.POSTRequestRawBody, headers, params, response)
	if err != nil || statusCode != 200 {
		logging.Error(fmt.Sprintf("Error in publishing data on pandemonium: %v | response %v", err, string(respBytes)))
	}
}
