package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func MemberChatroom(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the api/fetch_user_chatrooms request
	requestParams := map[string]string{
		ParamUserId: c.Query(ParamUserId),
		ParamUUID:   c.Query(ParamUUID),
		ParamState:  c.Query(ParamState),
		ParamPage:   c.Query(ParamPage),
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchMemberChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode)

}
