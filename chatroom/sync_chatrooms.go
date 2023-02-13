package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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
		ParamPage:          c.Query(ParamPage),
		ParamPageSize:      c.Query(ParamPageSize),
		ParamMaxTimeStamp:  c.Query(ParamMaxTimeStamp),
		ParamMinTimeStamp:  c.Query(ParamMinTimeStamp),
		ParamChatroomTypes: c.Query(ParamChatroomTypes),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, SyncChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
