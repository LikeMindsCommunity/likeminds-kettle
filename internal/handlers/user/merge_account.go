package user

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// MergeAccount used when user wants to merge account and generate login and refresh tokens
func MergeAccount(c *gin.Context) {

	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the merge account api internally
	body := map[string]string{}

	//Send Request
	utils.SendRequest(c, utils.CoreService, MergeAccountEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, body)
}
