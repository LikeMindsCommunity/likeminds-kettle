package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GetReportTags is used to fetch report tags in a community
func GetReportTags(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the member state api internally
	params := map[string]string{
		ParamType: c.Query(ParamType),
	}

	//Params Validation
	if params[ParamType] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchReportTagsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
