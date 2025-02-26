package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubCommon"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubPublish"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// SyncChatrooms is used to fetch data for chatroom syncing
func SyncChatrooms(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the sync chatroom api internally
	params := map[string]string{
		ParamPage:                       c.Query(ParamPage),
		ParamPageSize:                   c.Query(ParamPageSize),
		ParamMaxTimeStamp:               c.Query(ParamMaxTimeStamp),
		ParamMinTimeStamp:               c.Query(ParamMinTimeStamp),
		ParamChatroomTypes:              c.Query(ParamChatroomTypes),
		ParamIsLocalDb:                  c.Query(ParamIsLocalDb),
		ParamIncludedConversationStates: c.Query(ParamIncludedConversationStates),
		ParamChatroomId:                 c.Query(ParamChatroomId),
		ParamTag:                        c.Query(ParamTag),
	}

	//Get Request response
	syncChatroomsAPIResponseBytes, syncChatroomsAPIResponseBytesStatusCode := utils.GetRequestResponse(c, utils.CoreService, SyncChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if syncChatroomsAPIResponseBytes == nil {
		return
	}

	//Parse and generate response
	syncChatroomsAPIResponse := utils.ValidateClientResponse(c, syncChatroomsAPIResponseBytes, syncChatroomsAPIResponseBytesStatusCode)
	if syncChatroomsAPIResponse != nil {
		utils.GenerateResponse(c, syncChatroomsAPIResponse.Response, true)
		if syncChatroomsAPIResponse.Success == true {
			headers := utils.CreateHeadersFromToken(c, userId, user.GetRequestingUserDeviceId(c))

			minTimeStamp := c.Query(ParamMinTimeStamp)
			maxTimeStamp := c.Query(ParamMaxTimeStamp)

			utils.SafeGo(func() {
				parseAndPublishDROnTopicTypeChatroom(headers, minTimeStamp, maxTimeStamp, syncChatroomsAPIResponse.Response)
			})
		}
	}
}

// parseAndPublishDROnTopicTypeChatroom to publish DR on TopicTypeChatroom
func parseAndPublishDROnTopicTypeChatroom(headers map[string]interface{}, minTimeStamp, maxTimeStamp string, syncChatroomsAPIResponseMap map[string]interface{}) {
	// After returning the syncChatroomsAPIResponseMap, run a loop around "id" present in "chatroomsData:[]"
	chatroomsData, ok := syncChatroomsAPIResponseMap["chatrooms_data"].([]interface{})
	if !ok || chatroomsData == nil {
		logging.Error("chatrooms_data key is missing or is not a valid slice in the syncChatroomsAPIResponseMap")
		return // Exit the function or handle the error as appropriate
	}

	for _, chatroomData := range chatroomsData {
		chatroomMap, ok := chatroomData.(map[string]interface{})
		if !ok || chatroomMap == nil {
			logging.Error("chatroomData item is not a valid map in chatrooms_data")
			continue // Skip this item and continue with the next
		}

		chatroomID := chatroomMap["id"]
		communityID := chatroomMap["community_id"]
		if chatroomID == nil {
			logging.Error("parseAndPublishDROnTopicTypeChatroom: chatroomData ID is missing")
			return
		}
		if communityID == nil {
			logging.Error("parseAndPublishDROnTopicTypeChatroom: community ID is missing")
			return
		}

		// Call PublishDROnTopicTypeChatroom for each "chatroomID"
		pubsubPublish.PublishDROnTopicTypeChatroom(pubsubCommon.TopicMessageDeliveredDR, headers, minTimeStamp, maxTimeStamp, chatroomID, communityID)
	}
}
