package feed

import (
	"fmt"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// FetchGroupFeed is used to fetch group feed by a user
func FetchGroupFeed(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent in the /feed/group request
	params := map[string]string{
		ParamPage:       c.Query(ParamPage),
		ParamPageSize:   c.Query(ParamPageSize),
		ParamFeedroomId: c.Query(ParamFeedroomId),
		ParamTopicIds:   c.Query(ParamTopicIds),
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

	//Param updatiion
	params[ParamUserIsCm] = fmt.Sprint(response.IsCm)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchGroupFeedEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["posts"]; ok {
		posts := value.([]interface{})

		redisClient := utils.GetRedisClientFromContext(c)

		user_data, userUniqueIds, err := utils.GetUsersMetaFromFeedData(redisClient, headers, posts, dataResponse)
		if err != nil {
			utils.GenerateResponse(c, nil, false)
			return
		}

		// Update user data in dataResponse
		dataResponse["users"] = user_data

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
