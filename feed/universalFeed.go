package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchUniversalFeedEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["posts"]; ok {
		posts := value.([]interface{})

		if value, ok := dataResponse["filtered_comments"]; ok {
			if commentData, ok := value.(map[string]interface{}); ok {

				for _, val := range commentData {
					posts = append(posts, val)
				}
			}
		}

		user_data, err := user.GetUsersMetaFromFeedData(utils.GetRedisClientFromContext(c), utils.CreateHeaders(c, userId), posts, dataResponse)
		if err != nil {
			utils.GenerateResponse(c, nil, false)
			return
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
