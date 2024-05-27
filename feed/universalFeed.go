package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
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

	headers := utils.CreateHeaders(c, userId)

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
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchUniversalFeedEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse, err := populateDataResponseForFeed(headers, utils.GetRedisClientFromContext(c), apiCR.Response)
	if err != nil {
		utils.GenerateResponse(c, nil, false)
		return
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

func populateDataResponseForFeed(headers map[string]interface{}, redisClient *redis.Client, dataResponse map[string]interface{},
) (map[string]interface{}, error) {

	if value, ok := dataResponse["posts"]; ok {

		posts := value.([]interface{})

		if value, ok := dataResponse["filtered_comments"]; ok {
			if commentData, ok := value.(map[string]interface{}); ok {
				for _, val := range commentData {
					posts = append(posts, val)
				}
			}
		}

		userData, userUniqueIds, err := utils.GetUsersMetaFromFeedData(redisClient, headers, posts, dataResponse)
		if err != nil {
			return dataResponse, err
		}

		//Update user data in dataResponse
		dataResponse["users"] = userData

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	return dataResponse, nil
}
