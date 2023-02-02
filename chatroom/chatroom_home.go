package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Fetch chatroom home
func GetChatroomHome(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the collabcard seen api internally
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchChatroomHomeEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

}
