package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RenameChatroomRequest struct {
	ChatroomID      interface{} `json:"chatroom_id" binding:"required"`
	Header          string      `json:"header"`
	FirstTimeRename bool        `json:"first_time_rename"`
}

// RenameChatroom is used to rename an existing chatroom
func RenameChatroom(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the api/chatroom_rename POST request
	renameChatroomRequest, err := parseRenameChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, RenameChatroomEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, renameChatroomRequest)
}

func parseRenameChatroomRequest(c *gin.Context) (*RenameChatroomRequest, error) {
	//POST body params
	var rcr RenameChatroomRequest

	if err := c.ShouldBindJSON(&rcr); err != nil {
		return nil, err
	}

	if rcr.ChatroomID != nil {
		rcr.ChatroomID = utils.ParseInterfaceToString(rcr.ChatroomID)
	}

	return &rcr, nil
}
