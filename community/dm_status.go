package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func DMStatus(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the api/community_member/can_dm request
	requestParams := map[string]string{
		RequestFromParam: c.Query(RequestFromParam),
		ParamMemberId:    c.Query(ParamMemberId),
		ChatroomIDParam:  c.Query(ChatroomIDParam),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, UserCanDMEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
