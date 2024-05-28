package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func ListDMChatrooms(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the api/community_member/fetch_dm_chatrooms request
	requestParams := map[string]string{
		ParamPage: c.Query(ParamPage),
		ParamTag:  c.Query(ParamTag),
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchDMChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)
}
