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
	communitiesMeta, ok := response["community_meta"].(map[string]interface{})
	if !ok || communitiesMeta == nil {
		logging.Error("parseAndPublishDROnTopicTypeChatroom: community_meta key is missing or is not a valid map in the response")
		return // Exit the function or handle the error as appropriate
	}
	var communityID interface{}
	for _, communityMeta := range communitiesMeta {
		// Each `communityMeta` should be a map containing "id" and other details
		communityMetaMap, ok := communityMeta.(map[string]interface{})
		if !ok {
			logging.Error("parseAndPublishDROnTopicTypeChatroom: community_meta entry is not a valid map")
			return // Exit if the entry is not in the expected format
		}

		// Retrieve "id" from the first valid entry in "community_meta"
		communityID = communityMetaMap["id"]
		if communityID == nil {
			logging.Error("parseAndPublishDROnTopicTypeChatroom: id is missing in the community_meta entry")
			return // Exit if "id" is missing
		}
		break // Exit the loop after processing the first entry
	}

	// Ensure communityID is valid before proceeding
	if communityID == nil {
		logging.Debug("parseAndPublishDROnTopicTypeChatroom: no valid community ID found in community_meta")
		return
	}

	pubsubPublish.PublishDROnTopicTypeChatroom(headers, minTimeStamp, maxTimeStamp, chatroomID, communityID)
}
