package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type MuteChatroomRequest struct {
	ChatroomID interface{} `json:"chatroom_id" binding:"required"`
	Value      *bool       `json:"value" binding:"required"`
}

// MuteChatroom is used to mute a specifid chatroom
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
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, MuteChatroomEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, muteChatroomRequest)
}

func parseMuteChatroomRequest(c *gin.Context) (*MuteChatroomRequest, error) {
	//POST body params
	var mcr MuteChatroomRequest

	if err := c.ShouldBindJSON(&mcr); err != nil {
		return nil, err
	}

	if mcr.ChatroomID != nil {
		mcr.ChatroomID = utils.ParseInterfaceToString(mcr.ChatroomID)
	}

	return &mcr, nil
}
