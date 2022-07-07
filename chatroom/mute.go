package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type MuteChatroomRequest struct {
	ChatroomID int64 `json:"chatroom_id" binding:"required"`
	Value      bool  `json:"value" binding:"required"`
}

//MuteChatroom is used to mute a specifid chatroom
func MuteChatroom(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the mute chatroom api internally
	muteChatroomRequest, err := parseMuteChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, MuteChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, muteChatroomRequest)
}

func parseMuteChatroomRequest(c *gin.Context) (*MuteChatroomRequest, error) {
	//POST body params
	var mcr MuteChatroomRequest

	if err := c.ShouldBindJSON(&mcr); err != nil {
		return nil, err
	}

	return &mcr, nil
}
