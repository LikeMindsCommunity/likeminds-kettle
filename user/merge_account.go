package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

//MergeAccount used when user wants to merge account and generate login and refresh tokens
func MergeAccount(c *gin.Context) {

	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the merge account api internally
	body := map[string]string{
		//TODO - get mobile number and country code
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, MergeAccountEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, body)
}
