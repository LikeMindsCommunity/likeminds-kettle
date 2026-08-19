package feedroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/chatroom"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// FeedroomFollow is used to follow a specific feedroom
func FeedroomFollow(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the api/collabcard_follow request
	params := map[string]string{
		chatroom.ParamCollabcardId: c.Query(ParamFeedroomId),
		chatroom.ParamMemberId:     userId,
		chatroom.ParamValue:        c.Query(chatroom.ParamValue),
	}

	//Params Validation
	if params[chatroom.ParamCollabcardId] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.CollabcardFollowEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
