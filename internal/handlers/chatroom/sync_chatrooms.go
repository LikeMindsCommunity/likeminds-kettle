package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsub_publish"
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
	utils.ParseResponse(c, respBytes, statusCode, true)

	go pubsub_publish.PublishDROnTopicTypeChatroom(c, userId)
}
