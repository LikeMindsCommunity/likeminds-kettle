package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

func UserCanDM(c *gin.Context) {

	// Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the api/community_member/can_dm request
	requestParams := map[string]string{
		RequestFromParam: c.Query(RequestFromParam),
		MemberIDParam:    c.Query(MemberIDParam),
		ChatroomIDParam:  c.Query(ChatroomIDParam),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, UserCanDMEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
