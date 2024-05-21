package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateDMRequest struct {
	MemberID interface{} `json:"member_id"`
	UUID     string      `json:"uuid"`
	Tag      string      `json:"tag"`
}

func parseCreateDMRequest(c *gin.Context) (*CreateDMRequest, error) {
	//POST body params
	var cdmr CreateDMRequest
	if err := c.ShouldBindJSON(&cdmr); err != nil {
		return nil, err
	}

	if cdmr.MemberID != nil {
		cdmr.MemberID = utils.ParseInterfaceToString(cdmr.MemberID)
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
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CreateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createDMRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)
}
