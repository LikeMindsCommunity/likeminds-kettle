package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type InitiateDMRequest struct {
	ChatroomID       interface{} `json:"chatroom_id"`
	ChatRequestState int         `json:"chat_request_state"`
	Text             string      `json:"text"`
	MemberID         interface{} `json:"member_id"`
}

func parseInitiateDMRequest(c *gin.Context) (*InitiateDMRequest, error) {
	//POST body params
	var idmr InitiateDMRequest
	if err := c.ShouldBindJSON(&idmr); err != nil {
		return nil, err
	}

	if idmr.ChatroomID != nil {
		idmr.ChatroomID = utils.ParseInterfaceToString(idmr.ChatroomID)
	}

	if idmr.MemberID != nil {
		idmr.MemberID = utils.ParseInterfaceToString(idmr.MemberID)
	}

	return &idmr, nil
}

func InitiatingDMRequest(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the intiate dm request internally
	initiateDMRequest, err := parseInitiateDMRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, InitiateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, initiateDMRequest)
}
