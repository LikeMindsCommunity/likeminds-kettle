package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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
