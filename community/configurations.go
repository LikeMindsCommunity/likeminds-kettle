package community

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CommunityConfiguration struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Value       map[string]interface{} `json:"value"`
}

type CommunityConfigurationsResponse struct {
	Success                 bool                     `json:"success"`
	CommunityConfigurations []CommunityConfiguration `json:"community_configurations"`
}

// Expose method to fetch community configurations for a community
func GetCommunityConfigurations(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	params := map[string]string{
		ParamConfigurationTypes: c.Query(ParamConfigurationTypes),
	}

	//Send Request to api/community/configurations
	utils.SendRequest(c, utils.CoreService, FetchCommunityConfigurationsEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

}

// Exposed utility method to check if profile widgets are enabled for a community
func ProfileWidgetsEnabled(c *gin.Context, userId string) (bool, error) {

	if userId == "" {
		userId = user.GetRequestingUserId(c)
	}
	headers := utils.CreateHeaders(c, userId)

	apiKey := utils.GetApiKeyFromRequest(c)
	if apiKey == "" {
		return false, errors.New(utils.ErrorApiKeyNotFound)
	}

	redisClient := utils.GetRedisClientFromContext(c)
	if redisClient == nil {
		return false, errors.New(utils.ErrorRedisClientFailed)
	}

	profileMetaConfigurations, err := getProfileMetaCommunityConfigurationsFromCache(redisClient, headers, apiKey)
	if err != nil {
		return false, err
	}

	// Check if profile widgets are enabled
	if widgetsEnabled, ok := profileMetaConfigurations.Value[ConfigurationsProfileMetaWidgetsEnabled]; ok {

		return widgetsEnabled.(bool), nil
	}

	return false, nil
}

// utility method to fetch community configurations from Cache and if not, fetch from Core Service and updates cache
func getProfileMetaCommunityConfigurationsFromCache(redisClient *redis.Client, headers map[string]interface{}, apiKey string) (*CommunityConfiguration, error) {

	profileMetaConfigurations, exists, err := fetchProfileMetaConfigurationsfromCache(redisClient, apiKey)
	if err != nil {
		return nil, err
	}

	//If configurations are not present in cache, fetch from internal service and update cache
	if !exists {

		// Send request to internal service to fetch profile_metadata configurations
		profileMetaConfigurations, err = getCommunityConfigurationInternal(headers, CommunityConfigurationProfileMetadata)
		if err != nil || profileMetaConfigurations == nil {
			return nil, err
		}

		// Save data to cache
		setProfileMetaConfigurationsInCache(redisClient, apiKey, profileMetaConfigurations)

	}

	return profileMetaConfigurations, nil
}

// Internal method to fetch a community configuration without context
func getCommunityConfigurationInternal(headers map[string]interface{}, configurationType string) (*CommunityConfiguration, error) {

	params := map[string]string{
		ParamConfigurationTypes: utils.ParseStringArrayToString([]string{configurationType}),
	}

	// Send request to internal service
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, FetchCommunityConfigurationsEndpoint, utils.GETRequest, headers, params, nil)
	if err != nil || statusCode != http.StatusOK {
		return nil, err
	}

	// Parse response
	var ccr CommunityConfigurationsResponse
	if err = json.Unmarshal(respBytes, &ccr); err != nil {
		return nil, err
	}

	if len(ccr.CommunityConfigurations) == 0 {
		return nil, errors.New(utils.ErrorCommunityConfigurationsNotFound)
	}

	communityConfiguration := ccr.CommunityConfigurations[0]

	return &communityConfiguration, nil
}

// utility method to fetch profile_meta configurations from Cache
func fetchProfileMetaConfigurationsfromCache(redisClient *redis.Client, apiKey string) (*CommunityConfiguration, bool, error) {

	cacheKey := fmt.Sprintf(cache.ProfileMetaConfigurationsCacheKey, apiKey)

	//Fetch profile_meta configurations from cache
	profileMetaValue, exists, err := cache.Get(redisClient, cacheKey)
	if !exists {
		return nil, exists, err
	}

	var profileMetaConfigurations CommunityConfiguration
	err = json.Unmarshal([]byte(profileMetaValue), &profileMetaConfigurations)
	if err != nil {
		return nil, exists, err
	}

	return &profileMetaConfigurations, exists, nil
}

// utility method to save profile_meta configurations in Cache
func setProfileMetaConfigurationsInCache(redisClient *redis.Client, apiKey string, profileMetaConfigurations *CommunityConfiguration) error {

	cacheKey := fmt.Sprintf(cache.ProfileMetaConfigurationsCacheKey, apiKey)

	//Save profile_meta configurations in cache
	parsedProfileMeta, err := json.Marshal(profileMetaConfigurations)
	if err != nil {
		return err
	}

	err = cache.Set(redisClient, cacheKey, parsedProfileMeta, cache.ProfileMetaConfigurationsCacheTTL*time.Hour)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("Saved profile_meta configurations in cache for api-key: %s", apiKey))

	return nil
}
