package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ActionPendingChatroomRequest struct {
	ChatroomID int64 `json:"chatroom_id" binding:"required"`
	Value      *bool `json:"value" binding:"required"`
	PreApprove bool  `json:"pre_approve"`
}

//ActionPendingChatroom is used to change the state of pending chatroom
func ActionPendingChatroom(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the api/action_pending_chatroom POST request
	actionPendingChatroomRequest, err := parseActionPendingChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ActionPendingChatroomEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, actionPendingChatroomRequest)
}

func parseActionPendingChatroomRequest(c *gin.Context) (*ActionPendingChatroomRequest, error) {
	//POST body params
	var apcr ActionPendingChatroomRequest

	if err := c.ShouldBindJSON(&apcr); err != nil {
		return nil, err
	}

	return &apcr, nil
}
