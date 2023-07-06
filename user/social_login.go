package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

func UserSocialLogin(c *gin.Context) {
	// Params to be sent in the api/user/social/login GET request
	params := map[string]string{
		ParamLoginType:        c.Query(ParamLoginType),
		ParamSocialLoginToken: c.Query(ParamSocialLoginToken),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, UserSocialLoginEndpoint, utils.GETRequest, utils.CreateHeaders(c, ""), params, nil)
}
