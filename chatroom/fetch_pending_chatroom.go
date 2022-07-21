package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchPendingChatroom is used to fetch pending chatrooms
func FetchPendingChatroom(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch pending chatroom api internally
	params := map[string]string{
		ParamCommunityId: c.Query(ParamCommunityId),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchPendingChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
