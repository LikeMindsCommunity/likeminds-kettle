package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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
		ChatroomIDParam: c.Query(ChatroomIDParam),
	}

	//Params Validation
	if params[chatroom.ParamChatroomId] == "" {
		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchMembersMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
	} else {
		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.GetTaggingListEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}
