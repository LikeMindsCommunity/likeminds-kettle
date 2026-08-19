package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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
