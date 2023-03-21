package search

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/conversation"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// MessageSearch is used to perform search on the messages
func MessageSearch(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the conversation search api internally
	params := map[string]string{
		ParamSearch:                    c.Query(ParamSearch),
		ParamPage:                      c.Query(ParamPage),
		ParamPageSize:                  c.Query(ParamPageSize),
		conversation.ParamFollowStatus: c.Query(conversation.ParamFollowStatus),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, conversation.ConversationSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
