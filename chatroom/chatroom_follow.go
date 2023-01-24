package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// ChatroomFollow is used to follow a specific chatroom
func ChatroomFollow(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the api/collabcard_follow request
	params := map[string]string{
		ParamCollabcardId: c.Query(ParamCollabcardId),
		ParamMemberId:     c.Query(ParamMemberId),
		ParamValue:        c.Query(ParamValue),
	}

	//Params Validation
	if params[ParamCollabcardId] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CollabcardFollowEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
