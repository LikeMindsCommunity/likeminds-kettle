package utility

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// DecodeUrl is used to decode og tags for a url
func DecodeUrl(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the api/decode_url request
	params := map[string]string{
		ParamUrl: c.Query(ParamUrl),
	}

	//Params Validation
	if params[ParamUrl] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, DecodeUrlEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
