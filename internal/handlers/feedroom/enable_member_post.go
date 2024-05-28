package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// EnableMemberPostRequest | member message setting schema
type EnableMemberPostRequest struct {
	FeedroomId interface{} `json:"feedroom_id" binding:"required"`
	Value      *bool       `json:"value" binding:"required"`
}

// EnableMemberPost is used to enable member post settings in feedroom
func EnableMemberPost(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the enable member post api internally
	enableMemberPostRequest, err := parseEnableMemberPostRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	enableMemberMessageRequest := chatroom.EnableMemberMessageRequest{
		ChatroomId: enableMemberPostRequest.FeedroomId,
		Value:      enableMemberPostRequest.Value,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.EnableMemberMessageEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, enableMemberMessageRequest)
}

func parseEnableMemberPostRequest(c *gin.Context) (*EnableMemberPostRequest, error) {
	//POST body params
	var emmr EnableMemberPostRequest

	if err := c.ShouldBindJSON(&emmr); err != nil {
		return nil, err
	}

	if emmr.FeedroomId != nil {
		emmr.FeedroomId = utils.ParseInterfaceToString(emmr.FeedroomId)
	}

	return &emmr, nil
}
