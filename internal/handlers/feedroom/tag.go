package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// GetTaggingList is used to fetch the tag members list for a specific feedroom
func GetTaggingList(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the get tag list api internally
	params := map[string]string{
		chatroom.ParamChatroomId: c.Query(ParamFeedroomId),
	}

	// Access URl param, and if exists update query params
	feedroom_id := c.Param((ParamFeedroomId))
	if feedroom_id != "" {
		params[chatroom.ParamChatroomId] = feedroom_id
	}

	//Params Validation
	if params[chatroom.ParamChatroomId] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.GetTaggingListEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
