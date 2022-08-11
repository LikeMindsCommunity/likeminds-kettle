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

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchPendingChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
