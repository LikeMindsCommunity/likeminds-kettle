package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// GetManagementTools is used to get community management tools
func GetManagementTools(c *gin.Context) {
	ManagementTool(c, utils.GETMethod)
}

// ManagementTool method handles management tools in a community
func ManagementTool(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchManagementToolsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	}
}
