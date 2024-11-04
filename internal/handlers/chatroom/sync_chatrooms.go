package chatroom

import (
	"github.com/gin-gonic/gin"
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
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, SyncChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
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

			go parseAndPublishDROnTopicTypeChatroom(headers, minTimeStamp, maxTimeStamp, apiCR.Response)
		}
	}
}

// parseAndPublishDROnTopicTypeChatroom to publish DR on TopicTypeChatroom
func parseAndPublishDROnTopicTypeChatroom(headers map[string]interface{}, minTimeStamp, maxTimeStamp string, response map[string]interface{}) {
	// After returning the response, run a loop around "id" present in "chatrooms:[]"
	chatrooms, ok := response["chatrooms_data"].([]interface{})
	if !ok {
		logging.Error("chatrooms_data key is missing or is not a valid slice in the response")
		return // Exit the function or handle the error as appropriate
	}

	for _, chatroom := range chatrooms {
		chatroomMap, ok := chatroom.(map[string]interface{})
		if !ok {
			logging.Error("chatroom item is not a valid map in chatrooms_data")
			continue // Skip this item and continue with the next
		}

		chatroomID := chatroomMap["id"]
		communityID := chatroomMap["community_id"]

		// Call PublishDROnTopicTypeChatroom for each "chatroomID"
		pubsubPublish.PublishDROnTopicTypeChatroom(headers, minTimeStamp, maxTimeStamp, chatroomID, communityID)
	}

}
