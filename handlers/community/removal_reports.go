package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetRemovalReports(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Get bot if platform type is dashboard
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send Request to api/community/remove/reports
	utils.SendRequest(c, utils.CoreService, FetchCommunityRemovalReports, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
