package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetUserFeedMeta(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	// Access query params and url generation
	userID := c.Param("user_id")

	if userID == "" {
		utils.GeneralBadRequestError(c, utils.ErrorInvalidUserId)
	}

	//Get user_unique_id from user_id internally
	userUUID, err := utility.GetUUIDInternally(headers, userID)

	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	endpoint := fmt.Sprintf(FetchUserFeedMetaEndPoint, userUUID)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endpoint, utils.GETRequest, headers, nil, nil)

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	// If flow succeeds
	dataResponse := apiCR.Response

	//Fetch user data for given user_unique_ids
	userUniqueIds := []string{userUUID}
	redisClient := utils.GetRedisClientFromContext(c)

	user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, userUniqueIds)
	if err != nil {
		utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
		return
	}

	dataResponse["users"] = user_data

	// if user Topics connection is enabled, fetch and update related data
	if utils.UserTopicsConnectionEnabled(redisClient, headers) {
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
