package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/logging"
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
		logging.Info(fmt.Sprintf("community_id not found in cache for api-key: %s", apiKey))
		return communityId
	}

	communityId, err = strconv.Atoi(string(value))
	if err != nil {
		logging.Error(fmt.Sprintf("error parsing community_id from cache %v", err))
	}

	return communityId
}

func saveCommunityIdAgainstApiKeyToCache(redisClient *redis.Client, apiKey string, communityId int) error {

	cacheKey := fmt.Sprintf(cache.CommunityIdAgainstApiKeyCacheKey, apiKey)
	communityIdString := strconv.Itoa(communityId)

	if err := cache.Set(redisClient, cacheKey, []byte(communityIdString), cache.CommunityIdAgainstApiKeyCacheTTL*time.Hour); err != nil {
		logging.Error(fmt.Sprintf("error saving community_id in cache for api-key: %s | err: %v", apiKey, err))
		return err
	}

	return nil
}

// Exposed Method to fetch Community ID from API Key (from cache if present else from API)
func FetchCommunityIdFromApiKey(redisClient *redis.Client, apiKey string) (int, error) {

	defer Timer("FetchCommunityIdFromApiKey")()

	// Fetch community_id from cache
	communityId := getCommunityIdAgainstApiKeyFromCache(redisClient, apiKey)
	if communityId == 0 {

		headers := map[string]interface{}{
			HeadersApiKey: apiKey,
		}

		// Fetch community_id from API
		response, err := AuthenticateAPIKeyInternally(headers, apiKey)
		if err != nil {
			return communityId, err
		}

		communityId, err = strconv.Atoi(response[ParamCommunityID])
		if err != nil {
			return communityId, err
		}

		// Save community_id against apiKey in cache
		SafeGo(func() { saveCommunityIdAgainstApiKeyToCache(redisClient, apiKey, communityId) })
	}

	return communityId, nil
}
