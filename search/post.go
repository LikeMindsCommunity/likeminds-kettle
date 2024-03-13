package search

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// PostSearch is used to perform search on the posts
func PostSearch(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, feed.VIEW_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request to get excluded chatrooms list on Caravan Service
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, community.CommunityExcludedChatroomsEndPoint, utils.GETRequest, headers, nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	excludedChatroomIds := []int{}

	chatroomIds, ok := dataResponse["chatroom_ids"]
	if ok {
		for _, chatroomId := range chatroomIds.([]interface{}) {
			excludedChatroomIds = append(excludedChatroomIds, int(chatroomId.(float64)))
		}
	}

	temp_params, _ := json.Marshal(excludedChatroomIds)

	//Params to be sent in the post search api internally
	params := map[string]string{
		ParamSearch:              c.Query(ParamSearch),
		ParamSearchType:          c.Query(ParamSearchType),
		ParamPage:                c.Query(ParamPage),
		ParamPageSize:            c.Query(ParamPageSize),
		ParamExcludedChatroomIds: fmt.Sprintf("%v", string(temp_params)),
		ParamUserIsCm:            fmt.Sprint(response.IsCm),
	}

	//Send Request
	respBytes, statusCode = utils.GetRequestResponse(c, utils.SwarmService, PostSearchEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR = utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse = apiCR.Response
	if value, ok := dataResponse["posts"]; ok {
		posts := value.([]interface{})
		user_ids := []string{}

		//Fetch posts user id
		for _, post_data := range posts {
			if user_unique_id, ok := post_data.(map[string]interface{})["uuid"]; ok {
				user_ids = append(user_ids, user_unique_id.(string))
			}
		}

		user_ids = utils.AppendRepostPostUsersFromFeedDataResponse(dataResponse, user_ids)

		redisClient := utils.GetRedisClientFromContext(c)

		//Fetch user data for given user_unique_ids
		user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, user_ids)
		if err != nil {
			utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
			return
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data

		// if user Topics connection is enabled, fetch and update related data
		if utils.UserTopicsConnectionEnabled(redisClient, headers) {
			dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, user_ids)
		}
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, false)
}
