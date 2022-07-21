package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//GetTaggingList is used to fetch the tag members list for a specific chatroom
func GetTaggingList(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the get tag list api internally
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	//Params Validation
	if params[ParamChatroomId] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, GetTaggingListEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
