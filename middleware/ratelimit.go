package middleware

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/environment"
	"github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CommunityBillingMeta struct {
	TierType int `json:"tier_type"`
}

type TierDataType struct {
	MaxRequestLimitValue int    `json:"max_request_limit_value"`
	TTL                  int    `json:"ttl"`
	RateLimitKeyName     string `json:"rate_limit_key_name"`
	ErrorMessage         string `json:"error_message"`
	TierType             int    `json:"tier_type"`
	TierValueType        int    `json:"tier_value_type"`
}

type TierTypeApiResponse struct {
	Success bool           `json:"success"`
	Data    []TierDataType `json:"data"`
}

func RateLimitingMiddleware(redisClient *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Check if rate limiting is enabled
		isRateLimitingEnabled := environment.GoDotEnvVariable("IS_RATE_LIMITING_ENABLED")
		if isRateLimitingEnabled != "true" {
			c.Next()
			return
		}

		headers := utils.CreateHeaders(c, "")
		apiKey := headers[utils.HeadersApiKey].(string)

		// Get CommunityId from apiKey
		communityId, err := utils.FetchCommunityIdFromApiKey(redisClient, apiKey)
		if err != nil {
			logging.Error(err)
			return
		}

		// Get Community Billing data from cache
		communityBillingData, err := FetchCommunityBillingData(redisClient, communityId, headers)
		if err != nil {
			logging.Error(err)
			return
		}

		//Get TierData Function
		tierType := communityBillingData.TierType
		tierData, err := FetchTierData(redisClient, communityId, headers, tierType)
		if err != nil {
			logging.Error(err)
			return
		}
		isAllowed, err := checkTierDataForCommunityId(tierData, communityId, redisClient)
		// If rate limit is exceeded abort and return API error
		if !isAllowed {
			utils.RateLimitError(c, err.Error())
			return
		}
		// If any internal error occurs	log and continue
		if err != nil {
			logging.Error(err)
			return
		}

		c.Next()
	}
}

func checkTierDataForCommunityId(tierData []TierDataType, communityId int, redisClient *redis.Client) (bool, error) {
	for _, limitFactor := range tierData {
		// Extract all the required values from tierData
		rateLimitCurrentValueKey := limitFactor.RateLimitKeyName + fmt.Sprintf("_%d", communityId)
		rateLimitValue := limitFactor.MaxRequestLimitValue
		rateLimitErrorMessage := limitFactor.ErrorMessage
		rateLimitTTL := limitFactor.TTL
		rateLimitTierValueType := limitFactor.TierValueType

		// Get rate limit current value from cache
		currentValue, exists, err := cache.Get(redisClient, rateLimitCurrentValueKey)
		if err != nil {
			return true, err
		}

		// If key does not exist in cache
		if !exists {
			err = cache.Increment(redisClient, rateLimitCurrentValueKey)
			if err != nil {
				return true, err
			}
			redisClient.ExpireAt(rateLimitCurrentValueKey, time.Now().Add(time.Second*time.Duration(rateLimitTTL)))
			currentValue = "1"
		}

		// Check and calculate rate limit based on tier value type
		err = checkRateLimit(rateLimitTierValueType, currentValue, rateLimitValue, rateLimitErrorMessage, rateLimitCurrentValueKey, redisClient)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func checkRateLimit(rateLimitTierValueType int, currentValue string, rateLimitValue int, rateLimitErrorMessage string, rateLimitCurrentValueKey string, redisClient *redis.Client) error {
	currentValueInt, _ := strconv.Atoi(currentValue)
	switch rateLimitTierValueType {
	//RPM
	case RPM:
		// If rate limit current value is less than rate limit value
		if currentValueInt > rateLimitValue {
			return fmt.Errorf(rateLimitErrorMessage)
		}
		err := cache.Increment(redisClient, rateLimitCurrentValueKey)
		return err
	}
	return nil
}

func FetchCommunityBillingData(redisClient *redis.Client, communityId int, headers map[string]interface{}) (CommunityBillingMeta, error) {
	// Cache key for community billing data
	cacheKey := fmt.Sprintf(cache.CommunityBillingDataKey, communityId)
	// Get Community Billing data from cache
	value, valueExists, err := cache.Get(redisClient, cacheKey)
	//If error continue to next middleware
	if err != nil {
		return CommunityBillingMeta{}, err
	}

	communityBillingMeta := CommunityBillingMeta{}

	// communityBillingDataApi
	if !valueExists {
		// Get Value from API
		respBytes, _, err := utils.GetRequestResponseWithoutContext(utils.SubscriptionService, fmt.Sprintf("%s/%d", utils.BillingPlanEnpoint, communityId), utils.GETRequest, headers, nil, nil)
		if err != nil {
			return communityBillingMeta, err
		}

		err = json.Unmarshal(respBytes, &communityBillingMeta)
		if err != nil {
			return communityBillingMeta, err
		}

		communityBillingMetaForCache, err := json.Marshal(communityBillingMeta)

		if err != nil {
			return communityBillingMeta, err
		}
		// Update value in Cache
		err = cache.Set(redisClient, cacheKey, communityBillingMetaForCache, time.Hour*cache.CommunityBillingDataTTL)
		if err != nil {
			return communityBillingMeta, err
		}

	} else {
		// Unmarshal value from cache
		err := json.Unmarshal([]byte(value), &communityBillingMeta)
		if err != nil {
			return communityBillingMeta, err
		}
	}
	return communityBillingMeta, nil
}

func FetchTierData(redisClient *redis.Client, communityId int, headers map[string]interface{}, tierType int) ([]TierDataType, error) {

	cacheKey := fmt.Sprintf(cache.TierDataKey, tierType)
	value, exists, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	tierData := []TierDataType{}

	if !exists {
		params := map[string]string{
			utils.ParamTierType: strconv.Itoa(tierType),
		}
		// Get data from skulk service
		respBytes, _, err := utils.GetRequestResponseWithoutContext(utils.SubscriptionService, utils.TierEndpoint, utils.GETRequest, headers, params, nil)
		if err != nil {
			return nil, err
		}

		apiResponse := TierTypeApiResponse{}
		err = json.Unmarshal(respBytes, &apiResponse)
		if err != nil {
			return nil, err
		}

		cacheDataVal, err := json.Marshal(apiResponse.Data)
		if err != nil {
			return nil, err
		}

		// Save in cache
		err = cache.Set(redisClient, cacheKey, cacheDataVal, time.Hour*cache.TierDataTTL)
		if err != nil {
			return nil, err
		}
		tierData = apiResponse.Data
	} else {
		err := json.Unmarshal([]byte(value), &tierData)
		if err != nil {
			return nil, err
		}
	}
	return tierData, nil
}
