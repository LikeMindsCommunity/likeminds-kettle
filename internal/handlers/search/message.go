package search

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/conversation"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, conversation.ConversationSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}
