package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetIntroExamples(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, FetchIntroExamplesEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
