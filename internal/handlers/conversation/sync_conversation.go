package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubPublish"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// SyncConversation is used to sync conversations
func SyncConversation(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the sync conversation api internally
	params := map[string]string{
		ParamPage:                       c.Query(ParamPage),
		ParamPageSize:                   c.Query(ParamPageSize),
		ParamMaxTimeStamp:               c.Query(ParamMaxTimeStamp),
		ParamMinTimeStamp:               c.Query(ParamMinTimeStamp),
		ParamChatroomId:                 c.Query(ParamChatroomId),
		ParamIsLocalDb:                  c.Query(ParamIsLocalDb),
		ParamConversationId:             c.Query(ParamConversationId),
		ParamExcludedConversationStates: c.Query(ParamExcludedConversationStates),
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, SyncConversationEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR != nil {
		utils.GenerateResponse(c, apiCR.Response, true)
		if apiCR.Success == true {
			headers := utils.CreateHeadersFromToken(c, userId, user.GetRequestingUserDeviceId(c))
			minTimeStamp := c.Query(ParamMinTimeStamp)
			maxTimeStamp := c.Query(ParamMaxTimeStamp)
			chatroomID := c.Query(ParamChatroomId)

			go parseAndPublishDROnTopicTypeChatroom(headers, minTimeStamp, maxTimeStamp, chatroomID, apiCR.Response)
		}
	}
}

// parseAndPublishDROnTopicTypeChatroom to publish DR on TopicTypeChatroom
func parseAndPublishDROnTopicTypeChatroom(headers map[string]interface{}, minTimeStamp, maxTimeStamp, chatroomID string, response map[string]interface{}) {
	// Check if "community_meta" exists and is a map
	communityMeta, ok := response["community_meta"].(map[string]interface{})
	if !ok || communityMeta == nil {
		logging.Error("parseAndPublishDROnTopicTypeChatroom: community_meta key is missing or is not a valid map in the response")
		return // Exit the function or handle the error as appropriate
	}
	// Proceed to extract the community ID if community_meta is present
	communityID, ok := communityMeta["id"]
	if !ok || communityID == nil {
		logging.Error("parseAndPublishDROnTopicTypeChatroom: community ID is missing in community_meta")
		return // Handle missing community ID as appropriate
	}

	pubsubPublish.PublishDROnTopicTypeChatroom(headers, minTimeStamp, maxTimeStamp, chatroomID, communityID)
}
