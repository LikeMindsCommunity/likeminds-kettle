package pubsubPublish

import (
	"fmt"
	"github.com/nateshr/likeminds-authentication/internal/handlers/common"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// PublishDROnTopicTypeChatroom handles calling the Pandemonium publish API for inactive user flow in delivered report
func PublishDROnTopicTypeChatroom(publishAPITopicMessageType string, headers map[string]interface{}, minTimeStamp, maxTimeStamp string, chatroomID, communityID interface{}) {
	if chatroomID == nil {
		logging.Error("PublishDROnTopicTypeChatroom: chatroom ID is missing")
		return
	}
	if communityID == nil {
		logging.Error("PublishDROnTopicTypeChatroom: community ID is missing")
		return
	}

	topicChatroom := fmt.Sprintf(common.TopicTypeChatroomDynamic, chatroomID)

	publishAPIParams := map[string]string{
		common.ParamTopicMessageType: publishAPITopicMessageType,
	}

	// Prepare the publishAPIBodyMap (body) for the inactive user case
	publishAPIBodyMap := map[string]interface{}{
		common.ParamMinTimeStamp: minTimeStamp,
		common.ParamMaxTimeStamp: maxTimeStamp,
		common.ParamCommunityID:  communityID,
	}

	// Call the Pandemonium Publish API using your pre-defined function
	publishDataOnPandemonium(topicChatroom, headers, publishAPIParams, publishAPIBodyMap)
}

func publishDataOnPandemonium(topicChatroom string, headers map[string]interface{}, publishAPIParams map[string]string, publishAPIBodyMap map[string]interface{}) {
	publishAPIResponseBytes, publishAPIResponseStatusCode, err := utils.GetRequestResponseWithoutContext(utils.PandemoniumService, fmt.Sprintf(common.PublishEndPoint, topicChatroom), utils.POSTRequestRawBody, headers, publishAPIParams, publishAPIBodyMap)
	if err != nil || publishAPIResponseStatusCode != 200 {
		logging.Error(fmt.Sprintf(common.ErrorPublishAPI, err, string(publishAPIResponseBytes)))
	}
}
