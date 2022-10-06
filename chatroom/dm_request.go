package chatroom

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type InitiateDMRequest struct {
	ChatroomID       int    `json:"chatroom_id"`
	ChatRequestState int    `json:"chat_request_state"`
	Text             string `json:"text"`
	MemberID         int    `json:"member_id"`
}

func parseInitiateDMRequest(c *gin.Context) (*InitiateDMRequest, error) {
	//POST body params
	var idmr InitiateDMRequest
	if err := c.ShouldBindJSON(&idmr); err != nil {
		return nil, err
	}

	return &idmr, nil
}

func InitiatingDMRequest(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	fmt.Printf("%v", userId)

	//Params to be sent in the intiate dm request internally
	initiateDMRequest, err := parseInitiateDMRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, InitiateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, initiateDMRequest)
}
