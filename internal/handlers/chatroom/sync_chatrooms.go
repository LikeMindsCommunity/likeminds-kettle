package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubPublish"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
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
			go parseAndPublishDROnTopicTypeChatroom(c, userId, apiCR.Response)
		}
	}
}

// parseAndPublishDROnTopicTypeChatroom to publish DR on TopicTypeChatroom
func parseAndPublishDROnTopicTypeChatroom(c *gin.Context, userId string, response map[string]interface{}) {
	deviceID := user.GetRequestingUserDeviceId(c)

	// After returning the response, run a loop around "id" present in "chatrooms:[]"
	chatrooms := response["chatrooms_data"].([]interface{})

	for _, chatroom := range chatrooms {
		chatroomMap := chatroom.(map[string]interface{})
		chatroomID := chatroomMap["id"]
		communityID := chatroomMap["community_id"]

		// Call PublishDROnTopicTypeChatroom for each "chatroomID"
		pubsubPublish.PublishDROnTopicTypeChatroom(c, userId, deviceID, chatroomID, communityID)
	}

}
