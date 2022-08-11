package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchShareUrl is used to fetch share url for a specific chatroom
func FetchShareUrl(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch share url api internally
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchShareUrlEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
