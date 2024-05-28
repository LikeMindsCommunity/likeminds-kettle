package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// AutoFollowMembersRequest
type AutoFollowMembersRequest struct {
	ChatroomId          interface{} `json:"chatroom_id" binding:"required"`
	AutoFollowDone      *bool       `json:"auto_follow_done" binding:"required"`
	IncludeMembersLater *bool       `json:"include_members_later" binding:"required"`
}

// AutoFollowMembers is used to enable auto follow members for a chatroom
func AutoFollowMembers(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the auto follow for all members api internally
	autoFollowMembersRequest, err := parseAutoFollowMembersRequst(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, AutoFollowForAllMembersEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, autoFollowMembersRequest)
}

func parseAutoFollowMembersRequst(c *gin.Context) (*AutoFollowMembersRequest, error) {
	//POST body params
	var afmr AutoFollowMembersRequest

	if err := c.ShouldBindJSON(&afmr); err != nil {
		return nil, err
	}

	if afmr.ChatroomId != nil {
		afmr.ChatroomId = utils.ParseInterfaceToString(afmr.ChatroomId)
	}

	return &afmr, nil
}
