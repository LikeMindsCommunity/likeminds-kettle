package user

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/environment"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type MemberAccessResponse struct {
	Access bool `json:"access"`
	IsCm   bool `json:"is_cm"`
}

// FetchMemberAccess | fetch member access for sent action
// func FetchMemberAccess(c *gin.Context, accessType string, userId string) (bool, *MemberAccessResponse) {

// 	//Params to be sent in the api/community_member/fetch_access request
// 	params := map[string]string{
// 		ParamAccessType: accessType,
// 	}

// 	//Params Validation
// 	if params[ParamAccessType] == "" {
// 		//If GET params are missing
// 		utils.GETQueryParamsMissingError(c)
// 		return false, nil
// 	}

// 	//Send Request
// 	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchUserAccessEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

// 	//Validate response
// 	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
// 	if apiCR == nil {
// 		return false, nil
// 	}

// 	//If flow succeeds
// 	dataResponse := apiCR.Response
// 	response := MemberAccessResponse{
// 		Access: false,
// 		IsCm:   false,
// 	}

// 	//Create Data
// 	if value, ok := dataResponse["access"]; ok {
// 		response.Access = value.(bool)
// 	}

// 	if value, ok := dataResponse["is_cm"]; ok {
// 		response.IsCm = value.(bool)
// 	}

// 	return true, &response
// }

// FetchMemberAccess | fetch member access for sent action
func FetchMemberAccess(c *gin.Context, accessType string, userId string) (bool, *MemberAccessResponse) {

	cachedResponse := getAccessDataAgainstUserIdAndAccessTypeFromCache(utils.GetRedisClientFromContext(c), accessType, userId)
	if cachedResponse != nil {
		return true, cachedResponse
	}

	// Params for API request
	params := map[string]string{
		ParamAccessType: accessType,
	}

	// Params validation
	if params[ParamAccessType] == "" {
		utils.GETQueryParamsMissingError(c)
		return false, nil
	}

	// Send API request
	respBytes, statusCode := utils.GetRequestResponse(
		c, utils.CoreService, FetchUserAccessEndPoint, utils.GETRequest,
		utils.CreateHeaders(c, userId), params, nil,
	)

	// Validate API response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return false, nil
	}

	// Extract data from response
	dataResponse := apiCR.Response
	response := &MemberAccessResponse{
		Access: false,
		IsCm:   false,
	}

	if value, ok := dataResponse["access"]; ok {
		response.Access = value.(bool)
	}
	if value, ok := dataResponse["is_cm"]; ok {
		response.IsCm = value.(bool)
	}

	// Store response in Redis for 5 minutes
	setAccessDataAgainstUserIdAndAccessTypeFromCache(utils.GetRedisClientFromContext(c), accessType, userId, response)

	return true, response
}

func getAccessDataAgainstUserIdAndAccessTypeFromCache(redisClient *redis.Client, accessType string, userId string) *MemberAccessResponse {

	// Create a unique cache key for this request
	cacheKey := fmt.Sprintf(cache.FeedMemberAccessKey, userId, accessType)

	// Check Redis cache first
	cachedData, exists, err := cache.Get(redisClient, cacheKey)

	if err != nil {
		logging.Error(fmt.Sprintf("error fetching member access data from cache for key: %s, err: %v", cacheKey, err))
		return nil
	}

	if !exists {
		logging.Info(fmt.Sprintf("member access data not found in cache for userId_accessType: %s", cacheKey))
		return nil
	}
	// Parse JSON data
	var cachedResponse MemberAccessResponse
	if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err != nil {
		logging.Error(fmt.Sprintf("Error unmarshaling cached data for key: %s, err: %v", cacheKey, err))
		return nil
	}

	return &cachedResponse

}

func setAccessDataAgainstUserIdAndAccessTypeFromCache(redisClient *redis.Client, accessType string, userId string, data *MemberAccessResponse) error {
	// Create a unique cache key
	cacheKey := fmt.Sprintf(cache.FeedMemberAccessKey, userId, accessType)

	cacheExpirationTime := environment.GoDotEnvVariable("FEED_MEMBER_ACCESS_CACHE_TTL_IN_MIN")

	if cacheExpirationTime == "" {
		cacheExpirationTime = "5"
	}
	intCacheExpirationTime, err := strconv.Atoi(cacheExpirationTime)
	if err != nil {
		logging.Error(fmt.Sprintf("Error converting cache expiration time to int for key: %s, err: %v", cacheKey, err))
		return err
	}

	// Convert the response to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		logging.Error(fmt.Sprintf("Error marshaling access data for key: %s, err: %v", cacheKey, err))
		return err
	}

	// Store in Redis with expiration (e.g., 5 minutes)
	err = cache.Set(redisClient, cacheKey, string(jsonData), time.Minute*time.Duration(intCacheExpirationTime))
	if err != nil {
		logging.Error(fmt.Sprintf("Error setting access data in cache for key: %s, err: %v", cacheKey, err))
		return err
	}

	return nil
}

func DeleteAccessDataAgainstUserIdAndAccessTypeFromCache(redisClient *redis.Client, userId string) error {
	// Create a unique cache key
	cacheKey := fmt.Sprintf(cache.FeedMemberAccessKeyPattern, userId)

	// Use DeleteCacheByKeyPatterns to get the map of key patterns and matched keys.
	_, err := cache.DeleteCacheByKeyPatterns(redisClient, []string{cacheKey})
	if err != nil {
		logging.Error(fmt.Sprintf("Error deleting access data from cache for key: %s, err: %v", cacheKey, err))
		return err
	}
	logging.Info(fmt.Sprintf("Successfully deleted access data from cache for key: %s", cacheKey))
	return nil
}

func DeleteMultipleAccessDataAgainstUserIdAndAccessTypeFromCache(redisClient *redis.Client, userIds []string) error {
	var keyPatterns []string

	for _, userId := range userIds {
		keyPatterns = append(keyPatterns, fmt.Sprintf(cache.FeedMemberAccessKeyPattern, userId))
	}

	// Use DeleteCacheByKeyPatterns to get the map of key patterns and matched keys.
	_, err := cache.DeleteCacheByKeyPatterns(redisClient, keyPatterns)
	if err != nil {
		logging.Error(fmt.Sprintf("Error deleting access data from cache for multiple keys: %s, err: %v", keyPatterns, err))
		return err
	}
	logging.Info(fmt.Sprintf("Successfully deleted access data from cache for MULTIPLE keys: %s", keyPatterns))
	return nil
}

func DeleteAllAccessDataAgainstUserIdAndAccessTypeFromCache(redisClient *redis.Client) error {
	// Create a unique cache key
	cacheKey := cache.FeedMemberAccessKeyPrefix

	// Use DeleteCacheByKeyPatterns to get the map of key patterns and matched keys.
	_, err := cache.DeleteCacheByKeyPatterns(redisClient, []string{cacheKey})
	if err != nil {
		logging.Error(fmt.Sprintf("Error deleting access data from cache for key: %s, err: %v", cacheKey, err))
		return err
	}
	logging.Info(fmt.Sprintf("Successfully deleted access data from cache for ALL keys: %s", cacheKey))
	return nil
}
