package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PinChatroomRequest struct {
	ChatroomID int64 `json:"chatroom_id" binding:"required"`
	Value      *bool `json:"value" binding:"required"`
	Notify     bool  `json:"notify"`
}

//PinChatroom is used to create a pin a chatroom
func PinChatroom(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the pin chatroom request internally
	pinChatroomRequest, err := parsePinChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, PinChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, pinChatroomRequest)
}

func parsePinChatroomRequest(c *gin.Context) (*PinChatroomRequest, error) {
	//POST body params
	var pcr PinChatroomRequest

	if err := c.ShouldBindJSON(&pcr); err != nil {
		return nil, err
	}

	return &pcr, nil
}
