package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/logging"
)

// Exposed utility method to authenticate API Key internally
func AuthenticateAPIKeyInternally(headers map[string]interface{}, api_key string) (map[string]string, error) {

	var response map[string]string

	if api_key == "" {
		return response, nil
	}

	// Internally call /api/sdk/authenticate
	respBytes, statusCode, err := GetRequestResponseWithoutContext(CoreService, SDKAuthenticateEndPoint, GETRequest, headers, nil, nil)

	if err != nil {
		return nil, err
	}

	dataResponse := ValidateClientResponseWithoutContext(respBytes, statusCode, err)

	// Parse response
	if dataResponse != nil {
		var communityId string

		if dataResponse[ParamCommunityID] != nil {
			communityId = ParseInterfaceToString(dataResponse[ParamCommunityID])
		}

		response = map[string]string{
			ParamCommunityID: communityId,
		}
	}
	return response, nil
}

func getCommunityIdAgainstApiKeyFromCache(redisClient *redis.Client, apiKey string) int {

	communityId := 0

	cacheKey := fmt.Sprintf(cache.CommunityIdAgainstApiKeyCacheKey, apiKey)
	value, exists, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		logging.Error(fmt.Sprintf("error fetching community_id from cache for api-key: %s", apiKey))
		return communityId
	}

	if !exists {
		logging.Error(fmt.Sprintf("community_id not found in cache for api-key: %s", apiKey))
		return communityId
	}

	if err := json.Unmarshal([]byte(value), &communityId); err != nil {
		logging.Error(fmt.Sprintf("error unmarshalling community_id from cache for api-key: %s", apiKey))
	}

	return communityId
}

func saveCommunityIdAgainstApiKeyToCache(redisClient *redis.Client, apiKey string, communityId int) error {

	cacheKey := fmt.Sprintf(cache.CommunityIdAgainstApiKeyCacheKey, apiKey)
	communityIdString := strconv.Itoa(communityId)

	if err := cache.Set(redisClient, cacheKey, communityIdString, cache.CommunityIdAgainstApiKeyCacheTTL); err != nil {
		logging.Error(fmt.Sprintf("error saving community_id in cache for api-key: %s", apiKey))
		return err
	}

	return nil
}

// Exposed Method to fetch Community ID from API Key (from cache if present else from API)
func FetchCommunityIdFromApiKey(redisClient *redis.Client, apiKey string) (int, error) {

	// Fetch community_id from cache
	communityId := getCommunityIdAgainstApiKeyFromCache(redisClient, apiKey)
	if communityId == 0 {

		headers := map[string]interface{}{
			HeadersApiKey: apiKey,
		}

		// Fetch community_id from API
		response, err := AuthenticateAPIKeyInternally(headers, apiKey)
		if err != nil {
			logging.Error(fmt.Sprintf("error fetching community_id from API for api-key: %s", apiKey))
		}

		communityId, err = strconv.Atoi(response[ParamCommunityID])
		if err != nil {
			logging.Error(fmt.Sprintf("error converting community_id from response for api-key: %s", apiKey))
		}

		// Save community_id against apiKey in cache
		go saveCommunityIdAgainstApiKeyToCache(redisClient, apiKey, communityId)
	}

	return communityId, nil
}
