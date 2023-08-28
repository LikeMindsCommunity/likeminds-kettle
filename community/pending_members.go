package community

import (
	"github.com/gin-gonic/gin"

	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetPendingCommunityMembers(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Get bot id from context
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Send request to /api/pending_members
	utils.SendRequest(c, utils.CoreService, FetchPendingMembersEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
