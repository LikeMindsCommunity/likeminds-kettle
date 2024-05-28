package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type ChatroomBlockRequest struct {
	ChatroomID interface{} `json:"chatroom_id" binding:"required"`
	Status     int         `json:"status"`
}

func parseChatroomBlockRequest(c *gin.Context) (*ChatroomBlockRequest, error) {
	//POST body params
	var cbr ChatroomBlockRequest
	if err := c.ShouldBindJSON(&cbr); err != nil {
		return nil, err
	}

	if cbr.ChatroomID != nil {
		cbr.ChatroomID = utils.ParseInterfaceToString(cbr.ChatroomID)
	}

	return &cbr, nil
}

func ChatroomBlock(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the block chatroom request internally
	chatroomBlockRequest, err := parseChatroomBlockRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ChatroomBlockEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, chatroomBlockRequest)
}
