package feedroom

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetFeedroomMenu(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the /api/v2/fetch_chatroom request
	params := map[string]string{
		chatroom.ParamChatroomId: c.Query(ParamFeedroomId),
		chatroom.ParamApiType:    strconv.Itoa(chatroom.SdkApiType),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.FetchChatroomV2EndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
