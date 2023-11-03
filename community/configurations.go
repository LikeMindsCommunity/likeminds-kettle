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

// Internal method to fetch a community configuration without context
func getCommunityConfigurationInternal(headers map[string]interface{}, configurationType string) (*CommunityConfiguration, error) {

	params := map[string]string{
		ParamConfigurationTypes: fmt.Sprintf("[%v]", configurationType),
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

func fetchProfileMetaConfigurationsfromCache(redisClient *redis.Client, communityId int) (*CommunityConfiguration, error) {

	cacheKey := fmt.Sprintf(cache.ProfileMetaConfigurationsCacheKey, communityId)

	//Fetch profile_meta configurations from cache
	profileMetaValue, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		return nil, err
	}

	var profileMetaConfigurations CommunityConfiguration
	err = json.Unmarshal([]byte(profileMetaValue), &profileMetaConfigurations)
	if err != nil {
		return nil, err
	}

	return &profileMetaConfigurations, nil
}

func setProfileMetaConfigurationsInCache(redisClient *redis.Client, communityId int, profileMetaConfigurations *CommunityConfiguration) error {

	cacheKey := fmt.Sprintf(cache.ProfileMetaConfigurationsCacheKey, communityId)

	//Save profile_meta configurations in cache
	parsedProfileMeta, err := json.Marshal(profileMetaConfigurations)
	if err != nil {
		return err
	}

	err = cache.Set(redisClient, cacheKey, parsedProfileMeta, cache.ProfileMetaConfigurationsCacheTTL*time.Hour)
	if err != nil {
		return err
	}

	logging.Info(fmt.Sprintf("Saved the profile_meta configurations in cache for %d", communityId))

	return nil
}

// utility method to fetch community configurations from Cache and if not, fetch from Core Service and updates cache
func GetProfileMetaCommunityConfigurationsFromCache(redisClient *redis.Client, headers map[string]interface{}, communityId int) (*CommunityConfiguration, error) {

	profileMetaConfigurations, err := fetchProfileMetaConfigurationsfromCache(redisClient, communityId)
	if err != nil {
		return nil, err
	}

	//If configurations are not present in cache, fetch from internal service and update cache
	if profileMetaConfigurations.Type != CommunityConfigurationProfileMetadata {

		// Send request to internal service to fetch profile_metadata configurations
		profileMetaConfigurations, err = getCommunityConfigurationInternal(headers, CommunityConfigurationProfileMetadata)
		if err != nil || profileMetaConfigurations == nil {
			return nil, err
		}

		// Save data to cache
		setProfileMetaConfigurationsInCache(redisClient, communityId, profileMetaConfigurations)

	}

	return profileMetaConfigurations, nil
}

// utility method to check if profile widgets are enabled for a community
func IsProfileWidgetsEnabled(c *gin.Context) (bool, error) {

	headers := utils.CreateHeaders(c, user.GetRequestingUserId(c))

	communityId := utils.GetCommunityIdFromContext(c)
	if communityId == 0 {
		return false, errors.New(utils.ErrorCommunityIdNotFound)
	}

	redisClient := utils.GetRedisClientFromContext(c)
	if redisClient == nil {
		return false, errors.New(utils.ErrorRedisClientFailed)
	}

	profileMetaConfigurations, err := GetProfileMetaCommunityConfigurationsFromCache(redisClient, headers, communityId)
	if err != nil {
		return false, err
	}

	// Check if profile widgets are enabled
	if profileMetaConfigurations.Value[ConfigurationsProfileMetaWidgetsEnabled].(bool) {
		return true, nil
	}

	return false, nil
}
