package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//ChatroomSearch is used to perform search on the chatrooms
func ChatroomSearch(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Params to be sent in the chatroom search api internally
	params := map[string]string{
		ParamSearch:       c.Query(ParamSearch),
		ParamFollowStatus: c.Query(ParamFollowStatus),
		ParamPage:         c.Query(ParamPage),
		ParamPageSize:     c.Query(ParamPageSize),
		ParamSearchType:   c.Query(ParamSearchType),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ChatroomSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
