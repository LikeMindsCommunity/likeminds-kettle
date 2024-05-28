package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func GetIntroExamples(c *gin.Context) {
	var userId string

	ltm, _ := c.Get(constants.ParamLTM)

	if ltm != nil {
		// Authorize User
		userId := user.GetRequestingUserId(c)
		if userId == "" {
			return
		}
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, FetchIntroExamplesEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
