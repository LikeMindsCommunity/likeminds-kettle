package pubsubPublish

import (
	"fmt"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/logging"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
)

// PublishConversationOnTopicTypeChatroom to publish Conversation on TopicTypeChatroomDynamic
func PublishConversationOnTopicTypeChatroom(headers map[string]interface{}, chatroomID interface{}, response map[string]interface{}) {
	if chatroomID == nil {
		logging.Error("PublishConversationOnTopicTypeChatroom: chatroom ID is missing")
		return
	}
	if response == nil {
		logging.Error("PublishConversationOnTopicTypeChatroom: response is missing")
		return
	}

	topicChatroom := fmt.Sprintf(TopicTypeChatroomDynamic, chatroomID)
	params := map[string]string{
		ParamTopicMessageType: TopicMessageTypeConversation,
	}

	publishDataOnPandemonium(topicChatroom, headers, params, response)
}

// PublishConversationOnTopicTypeCommunity to publish Conversation on TopicTypeCommunityDynamic
func PublishConversationOnTopicTypeCommunity(headers map[string]interface{}, response map[string]interface{}) {
	if response == nil {
		logging.Error("PublishConversationOnTopicTypeCommunity: response is missing")
		return
	}

	// Check if "conversation" exists and is a map
	conversation, ok := response["conversation"].(map[string]interface{})
	if !ok {
		logging.Error("PublishConversationOnTopicTypeCommunity: conversation key is missing or is not a valid map in response")
		return
	}
	// Check if "member" exists within "conversation" and is a map
	member, ok := conversation["member"].(map[string]interface{})
	if !ok {
		logging.Error("PublishConversationOnTopicTypeCommunity: member key is missing or is not a valid map in conversation")
		return
	}
	// Check if "sdk_client_info" exists within "member" and is a map
	sdkClientInfo, ok := member["sdk_client_info"].(map[string]interface{})
	if !ok {
		logging.Error("PublishConversationOnTopicTypeCommunity: sdk_client_info key is missing or is not a valid map in member")
		return
	}
	// Check if "community" exists within "sdk_client_info" and is of the expected type
	communityID, ok := sdkClientInfo["community"].(float64)
	if !ok {
		logging.Error("PublishConversationOnTopicTypeCommunity: community key is missing or is not a valid string in sdk_client_info")
		return
	}

	topicCommunity := fmt.Sprintf(TopicTypeCommunityDynamic, communityID)
	publishParams := map[string]string{
		ParamTopicMessageType: TopicMessageTypeConversation,
	}

	publishDataOnPandemonium(topicCommunity, headers, publishParams, response)
}

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

	topicChatroom := fmt.Sprintf(TopicTypeChatroomDynamic, chatroomID)
	params := map[string]string{
		ParamTopicMessageType: TopicMessageDR,
	}

	// Prepare the response (body) for the inactive user case
	response := map[string]interface{}{
		ParamMinTimeStamp: minTimeStamp,
		ParamMaxTimeStamp: maxTimeStamp,
		ParamCommunityID:  communityID,
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
