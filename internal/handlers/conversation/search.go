package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// ConversationSearch is used to perform search on the conversations
func ConversationSearch(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the conversation search api internally
	params := map[string]string{
		ParamSearch:       c.Query(ParamSearch),
		ParamFollowStatus: c.Query(ParamFollowStatus),
		ParamPage:         c.Query(ParamPage),
		ParamPageSize:     c.Query(ParamPageSize),
		ParamChatroomId:   c.Query(ParamChatroomId),
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, ConversationSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}
