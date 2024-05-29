package middleware

import (
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
		communityBillingData, err := utils.FetchCommunityBillingData(redisClient, communityId, headers)
		if err != nil {
			logging.Error(err)
			return
		}

		//Get TierData Function
		tierType := communityBillingData.TierType
		tierData, err := utils.FetchTierData(redisClient, headers, tierType)
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

func checkTierDataForCommunityId(tierData []utils.TierDataType, communityId int, redisClient *redis.Client) (bool, error) {
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

			// Set expiration time for the key
			err = cache.ExpireAt(redisClient, rateLimitCurrentValueKey, time.Now().Add(time.Second*time.Duration(rateLimitTTL)))
			if err != nil {
				return true, err
			}

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
