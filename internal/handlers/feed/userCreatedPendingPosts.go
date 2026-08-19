package feed

import (
	"fmt"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/utility"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// FetchUserCreatedPendingPost is used to fetch pending posts created by a user
func FetchUserCreatedPendingPosts(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent in the /user/<user_id>/post request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Access query params and url generation
	paramUUID := c.Param("uuid")

	//Get user_unique_id from paramUUID internally
	uuid, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), paramUUID)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION, userId)
	if !success {
		utils.MemberAccessFailError(c)
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	if !response.IsCm && userId != uuid {
		utils.MemberAccessFailError(c)
		return
	}

	//Url generation
	userCreatedPendingPostsEndPoint := fmt.Sprintf(FetchUserCreatedPendingPostsEndPoint, uuid)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, userCreatedPendingPostsEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	dataResponse = populatePendingPostsDataResponse(c, dataResponse)

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

// Internal method to populate users data
func populatePendingPostsDataResponse(c *gin.Context, dataResponse map[string]interface{}) map[string]interface{} {
	if value, ok := dataResponse["posts"]; ok {
		posts := value.([]interface{})
		user_ids := []string{}

		//Fetch posts user id
		for _, post_data := range posts {
			if user_unique_id, ok := post_data.(map[string]interface{})["uuid"]; ok {
				user_ids = append(user_ids, user_unique_id.(string))
			}
		}

		userId := user.GetRequestingUserId(c)
		redisClient := utils.GetRedisClientFromContext(c)
		headers := utils.CreateHeaders(c, userId)

		//Fetch user data for given user_unique_ids
		user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, user_ids)
		if err != nil {
			utils.GeneralAPIError(c, utils.ErrorFetchingUserData)
			return nil
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, user_ids)
	}

	return dataResponse
}
