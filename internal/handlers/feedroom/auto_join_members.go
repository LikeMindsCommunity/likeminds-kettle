package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// AutoJoinMembersRequest
type AutoJoinMembersRequest struct {
	FeedroomId          interface{} `json:"feedroom_id" binding:"required"`
	AutoFollowDone      *bool       `json:"auto_follow_done" binding:"required"`
	IncludeMembersLater *bool       `json:"include_members_later" binding:"required"`
}

// AutoJoinMembers is used to enable auto join members for a feedroom
func AutoJoinMembers(c *gin.Context) {

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
	autoJoinMembersRequest, err := parseAutoJoinMembersRequst(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	autoFollowMembersRequest := chatroom.AutoFollowMembersRequest{
		ChatroomId:          autoJoinMembersRequest.FeedroomId,
		AutoFollowDone:      autoJoinMembersRequest.AutoFollowDone,
		IncludeMembersLater: autoJoinMembersRequest.IncludeMembersLater,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.AutoFollowForAllMembersEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, autoFollowMembersRequest)
}

func parseAutoJoinMembersRequst(c *gin.Context) (*AutoJoinMembersRequest, error) {
	//POST body params
	var afmr AutoJoinMembersRequest

	if err := c.ShouldBindJSON(&afmr); err != nil {
		return nil, err
	}

	if afmr.FeedroomId != nil {
		afmr.FeedroomId = utils.ParseInterfaceToString(afmr.FeedroomId)
	}

	return &afmr, nil
}
