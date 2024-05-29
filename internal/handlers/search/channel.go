package search

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// ChannelSearch is used to perform search on the channels
func ChannelSearch(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the chatroom search api internally
	params := map[string]string{
		ParamSearch:                c.Query(ParamSearch),
		ParamSearchType:            c.Query(ParamSearchType),
		ParamPage:                  c.Query(ParamPage),
		ParamPageSize:              c.Query(ParamPageSize),
		chatroom.ParamFollowStatus: c.Query(chatroom.ParamFollowStatus),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.ChatroomSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
