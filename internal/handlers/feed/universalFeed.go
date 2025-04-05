package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// FetchUniversalFeed is used to fetch universal feed by a user
func FetchUniversalFeed(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the /feed/universal request
	params := map[string]string{
		ParamPage:      c.Query(ParamPage),
		ParamPageSize:  c.Query(ParamPageSize),
		ParamTopicIds:  c.Query(ParamTopicIds),
		ParamWidgetIds: c.Query(ParamWidgetIds),
		ParamPostIds:   c.Query(ParamPostIds),
	}

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Param updation
	params[ParamUserIsCm] = fmt.Sprint(response.IsCm)

	//add CM role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	headers := utils.CreateHeaders(c, userId)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchUniversalFeedEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse, err := utils.PopulateDataResponseForFeed(headers, utils.GetRedisClientFromContext(c), apiCR.Response)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
