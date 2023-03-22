package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UserTokenRequest struct {
	Token string `json:"token"`
}

func PushUserToken(c *gin.Context) {
	// Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	var usrToken UserTokenRequest

	if err := c.ShouldBindJSON(&usrToken); err != nil {
		return
	}

	// Params to be sent in api/push
	params := map[string]string{
		ParamMemberId: userId,
		ParamDeviceId: c.GetHeader(utils.HeadersDeviceId),
		ParamToken:    usrToken.Token,
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, PushUserTokenEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
