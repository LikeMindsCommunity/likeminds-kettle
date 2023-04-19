package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// MyChatrooms is used to fetch all the chatrooms for a user
func MyChatrooms(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the my chatroom api internally
	params := map[string]string{
		ParamPage: c.Query(ParamPage),
		ParamTag:  c.Query(ParamTag),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, MyChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
