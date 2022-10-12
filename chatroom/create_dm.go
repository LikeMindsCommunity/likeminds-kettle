package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateDMRequest struct {
	MemberID int `json:"member_id"`
}

func parseCreateDMRequest(c *gin.Context) (*CreateDMRequest, error) {
	//POST body params
	var cdmr CreateDMRequest
	if err := c.ShouldBindJSON(&cdmr); err != nil {
		return nil, err
	}

	return &cdmr, nil
}

func CreateDM(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the create dm request internally
	createDMRequest, err := parseCreateDMRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CreateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createDMRequest)
}
