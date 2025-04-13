package feed

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// func GetUserFeedMeta(c *gin.Context) {
// 	// Authorize User
// 	userId := user.GetRequestingUserId(c)
// 	if userId == "" {
// 		return
// 	}

// 	headers := utils.CreateHeaders(c, userId)

// 	// Access query params and url generation
// 	paramUUID := c.Param("uuid")

// 	if paramUUID == "" {
// 		utils.GeneralBadRequestError(c, utils.ErrorInvalidUserId)
// 	}

// 	//Get user_unique_id from user_id internally
// 	userUUID, err := utility.GetUUIDInternally(headers, paramUUID)
// 	if err != nil {
// 		utils.GeneralAPIError(c, err.Error())
// 		return
// 	}

// 	endpoint := fmt.Sprintf(FetchUserFeedMetaEndPoint, userUUID)

// 	//Send Request
// 	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endpoint, utils.GETRequest, headers, nil, nil)

// 	// Validate response
// 	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
// 	if apiCR == nil {
// 		return
// 	}

// 	// If flow succeeds
// 	dataResponse := apiCR.Response

// 	//Fetch user data for given user_unique_ids
// 	userUniqueIds := []string{userUUID}
// 	redisClient := utils.GetRedisClientFromContext(c)

// 	user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, userUniqueIds)
// 	if err != nil {
// 		utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
// 		return
// 	}

// 	dataResponse["users"] = user_data

// 	// Update user topics data in dataResponse
// 	dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)

// 	//Send response
// 	utils.GenerateResponse(c, dataResponse, true)
// }

func GetUserFeedMeta(c *gin.Context) {

	timestamp1 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timestamp1 + " - timestamp1")

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	// Access query params and url generation
	paramUUID := c.Param("uuid")

	if paramUUID == "" {
		utils.GeneralBadRequestError(c, utils.ErrorInvalidUserId)
	}

	//Get user_unique_id from user_id internally
	userUUID, err := utility.GetUUIDInternally(headers, paramUUID)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	endpoint := fmt.Sprintf(FetchUserFeedMetaEndPoint, userUUID)

	timestamp2 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timestamp2 + " - timestamp2")

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endpoint, utils.GETRequest, headers, nil, nil)

	timestamp9 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timestamp9 + " - timestamp9")

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

	timestamp10 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timestamp10 + " - timestamp10")

	dataResponse["users"] = user_data

	// Update user topics data in dataResponse
	dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)

	timestamp11 := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(timestamp11 + " - timestamp11")

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
