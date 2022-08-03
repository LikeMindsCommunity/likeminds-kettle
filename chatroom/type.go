package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ChatroomTypeRequest struct {
	ChatroomID *int32 `json:"chatroom_id" binding:"required"`
	IsSecret   *bool  `json:"is_secret" binding:"required"`
}

//ChatroomType is used to change the type of chatroom
func ChatroomType(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the api/chatroom/change_type POST request
	chatroomTypeRequest, err := parseChatroomTypeRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ChatroomTypeEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, chatroomTypeRequest)
}

func parseChatroomTypeRequest(c *gin.Context) (*ChatroomTypeRequest, error) {
	//POST body params
	var ctr ChatroomTypeRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	return &ctr, nil
}
