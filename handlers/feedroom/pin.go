package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PinFeedroomRequest struct {
	FeedroomID interface{} `json:"feedroom_id" binding:"required"`
	Value      *bool       `json:"value" binding:"required"`
}

// PinFeedroom is used to create a pin a feedroom
func PinFeedroom(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the pin feedroom request internally
	pinFeedroomRequest, err := parsePinFeedroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	pinChatroomRequest := chatroom.PinChatroomRequest{
		ChatroomID: pinFeedroomRequest.FeedroomID,
		Value:      pinFeedroomRequest.Value,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.PinChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, pinChatroomRequest)
}

func parsePinFeedroomRequest(c *gin.Context) (*PinFeedroomRequest, error) {
	//POST body params
	var pcr PinFeedroomRequest

	if err := c.ShouldBindJSON(&pcr); err != nil {
		return nil, err
	}

	if pcr.FeedroomID != nil {
		pcr.FeedroomID = utils.ParseInterfaceToString(pcr.FeedroomID)
	}

	return &pcr, nil
}
