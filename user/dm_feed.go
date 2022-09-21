package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

func DMFeed(c *gin.Context) {

	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, UserDMFeedEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
