package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GetReportTags is used to fetch report tags in a community
func GetReportTags(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	apiRevampCheckv1 := utils.ApiRevampV1Check(c)

	// Query Params
	params := map[string]string{}

	if apiRevampCheckv1 {

		params[ParamEntityType] = c.Query(ParamEntityType)

		// Validation for entity_type
		if params[ParamEntityType] == "" {
			utils.GeneralBadRequestError(c, "Please provide entity_type in query params")
			return
		}

		//Send Request to api/community/report/tag with GET method
		utils.SendRequest(c, utils.CoreService, FetchCommunityReportTagsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {

		params[ParamType] = c.Query(ParamType)

		// Validation for type
		if params[ParamType] == "" {
			utils.GeneralBadRequestError(c, "Please provide type in query params")
			return
		}

		//Send Request to api/fetch_report_tags with GET method
		utils.SendRequest(c, utils.CoreService, FetchReportTagsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	}
}
