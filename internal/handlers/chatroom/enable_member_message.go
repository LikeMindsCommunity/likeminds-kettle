package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// EnableMemberMessageRequest | member message setting schema
type EnableMemberMessageRequest struct {
	ChatroomId interface{} `json:"chatroom_id" binding:"required"`
	Value      *bool       `json:"value" binding:"required"`
}

// EnableMemberMessage is used to enable member message settings in chatroom
func EnableMemberMessage(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the enable member message api internally
	enableMemberMessageRequest, err := parseEnableMemberMessageRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EnableMemberMessageEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, enableMemberMessageRequest)
}

func parseEnableMemberMessageRequest(c *gin.Context) (*EnableMemberMessageRequest, error) {
	//POST body params
	var emmr EnableMemberMessageRequest

	if err := c.ShouldBindJSON(&emmr); err != nil {
		return nil, err
	}

	if emmr.ChatroomId != nil {
		emmr.ChatroomId = utils.ParseInterfaceToString(emmr.ChatroomId)
	}

	return &emmr, nil
}
