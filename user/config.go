package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Config | fetch user app config
func Config(c *gin.Context) {

	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the api/config request
	params := map[string]string{
		ParamIngestCommunities: c.Query(ParamIngestCommunities),
	}

	//Params Validation
	if params[ParamIngestCommunities] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ConfigEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
