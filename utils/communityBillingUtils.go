package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/logging"
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

type BillingPlanApiResponse struct {
	Success     bool                 `json:"success"`
	BillingData CommunityBillingMeta `json:"billing_data"`
}

type TierTypeApiResponse struct {
	Success bool           `json:"success"`
	Data    []TierDataType `json:"data"`
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
		respBytes, _, err := GetRequestResponseWithoutContext(SubscriptionService, fmt.Sprintf("%s/%d", BillingPlanEnpoint, communityId), GETRequest, headers, nil, nil)
		if err != nil {
			return communityBillingMeta, err
		}

		billingPlanResp := BillingPlanApiResponse{}
		err = json.Unmarshal(respBytes, &billingPlanResp)
		communityBillingMeta = billingPlanResp.BillingData
		if err != nil {
			return communityBillingMeta, err
		}

		communityBillingMetaForCache, err := json.Marshal(billingPlanResp.BillingData)

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

func FetchTierData(redisClient *redis.Client, headers map[string]interface{}, tierType int) ([]TierDataType, error) {

	cacheKey := fmt.Sprintf(cache.TierDataKey, tierType)
	value, exists, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		logging.Error(err)
		return nil, err
	}

	tierData := []TierDataType{}

	if !exists {
		params := map[string]string{
			ParamTierType: strconv.Itoa(tierType),
		}
		// Get data from skulk service
		respBytes, _, err := GetRequestResponseWithoutContext(SubscriptionService, TierEndpoint, GETRequest, headers, params, nil)
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

// Check if community is on free tier | returns true if any error
func IsCommunityOnFreeTier(redisClient *redis.Client, headers map[string]interface{}) bool {

	// Get CommunityId from apiKey
	communityId, err := FetchCommunityIdFromApiKey(redisClient, headers[HeadersApiKey].(string))
	if err != nil {
		logging.Error(err)
		return true
	}

	// Get Community Billing meta from cache
	communityBillingMeta, err := FetchCommunityBillingData(redisClient, communityId, headers)
	if err != nil {
		logging.Error(err)
		return true
	}

	if communityBillingMeta.TierType == FreeTierType {
		return true
	} else {
		return false
	}
}
