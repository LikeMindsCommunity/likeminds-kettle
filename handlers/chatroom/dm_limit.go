package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func DMLimit(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the api/community_member/request_dm_limit request
	requestParams := map[string]string{
		ParamMemberId: c.Query(ParamMemberId),
		ParamUUID:     c.Query(ParamUUID),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, RequestDMLimitEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
