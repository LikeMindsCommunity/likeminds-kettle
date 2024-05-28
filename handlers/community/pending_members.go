package community

import (
	"github.com/gin-gonic/gin"

	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func GetPendingCommunityMembers(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request to /api/pending_members
	utils.SendRequest(c, utils.CoreService, FetchPendingMembersEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
