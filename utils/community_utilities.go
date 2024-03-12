package utils

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
)

type CommunityConfiguration struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Value       map[string]interface{} `json:"value"`
}

// CommunitySetting | schema for community settings
type CommunitySetting struct {
	SettingType     string `json:"setting_type"  binding:"required"`
	SettingTitle    string `json:"setting_title"  binding:"required"`
	SettingSubTitle string `json:"setting_sub_title"  binding:"required"`
	IsEnabled       bool   `json:"enabled"  binding:"required"`
}

type CommunitySettingsResponse struct {
	Success           bool               `json:"success"`
	CommunitySettings []CommunitySetting `json:"community_settings"`
}

type CommunityConfigurationsResponse struct {
	Success                 bool                     `json:"success"`
	CommunityConfigurations []CommunityConfiguration `json:"community_configurations"`
}

// Exposed utility method to check if profile widgets are enabled for a community
func IsProfileWidgetsEnabled(c *gin.Context, userId string) (bool, error) {

	headers := CreateHeaders(c, userId)

	apiKey := headers[HeadersApiKey].(string)
	if apiKey == "" {
		return false, errors.New(ErrorApiKeyNotFound)
	}

	redisClient := GetRedisClientFromContext(c)
	if redisClient == nil {
		return false, errors.New(ErrorRedisClientFailed)
	}

	profileMetaConfigurations, err := getProfileMetaConfig(redisClient, headers, apiKey)
	if err != nil || profileMetaConfigurations == nil {
		return false, err
	}

	// Check if profile widgets are enabled
	if widgetsEnabled, ok := profileMetaConfigurations.Value[ConfigurationsProfileMetaWidgetsEnabled]; ok {

		return widgetsEnabled.(bool), nil
	}

	return false, nil
}

// utility method to fetch community configurations from Cache and if not, fetch from Core Service and updates cache
func getProfileMetaConfig(redisClient *redis.Client, headers map[string]interface{}, apiKey string) (*CommunityConfiguration, error) {

	profileMetaConfig, exists, err := fetchProfileMetaConfigfromCache(redisClient, apiKey)
	if err != nil {
		return nil, err
	}

	//If configurations are not present in cache, fetch from internal service and update cache
	if !exists {

		// Send request to internal service to fetch profile_metadata configurations
		profileMetaConfig, err = getCommunityConfigurationInternal(headers, CommunityConfigurationProfileMetadata)
		if err != nil || profileMetaConfig == nil {
			return nil, err
		}

		// Save data to cache
		setProfileMetaConfigInCache(redisClient, apiKey, profileMetaConfig)

	}

	return profileMetaConfig, nil
}

// Internal method to fetch a community configuration without context
func getCommunityConfigurationInternal(headers map[string]interface{}, configurationType string) (*CommunityConfiguration, error) {

	params := map[string]string{
		ParamConfigurationTypes: ParseStringArrayToString([]string{configurationType}),
	}

	// Send request to internal service
	respBytes, statusCode, err := GetRequestResponseWithoutContext(CoreService, FetchCommunityConfigurationsEndpoint, GETRequest, headers, params, nil)
	if err != nil || statusCode != http.StatusOK {
		return nil, err
	}

	// Parse response
	var ccr CommunityConfigurationsResponse
	if err = json.Unmarshal(respBytes, &ccr); err != nil {
		return nil, err
	}

	if len(ccr.CommunityConfigurations) == 0 {
		return nil, errors.New(ErrorCommunityConfigurationsNotFound)
	}

	communityConfiguration := ccr.CommunityConfigurations[0]

	return &communityConfiguration, nil
}

// GetCommunitySettingsInternal | fetch community setting for application internal use
func GetCommunitySettingsInternal(headers map[string]interface{}) ([]CommunitySetting, error) {

	// Send request to internal service
	respBytes, statusCode, err := GetRequestResponseWithoutContext(CoreService, FetchCommunitySettingsEndpoint, GETRequest, headers, nil, nil)
	if err != nil || statusCode != http.StatusOK {
		return nil, err
	}

	// Parse response
	var csr CommunitySettingsResponse
	if err = json.Unmarshal(respBytes, &csr); err != nil {
		return nil, err
	}

	if len(csr.CommunitySettings) == 0 {
		return nil, errors.New(ErrorCommunitySettingsNotFound)
	}

	return csr.CommunitySettings, nil
}

// CheckCommunitySettingEnabled | check is setting type is enabled in the community
func CheckCommunitySettingEnabled(communitySettings []CommunitySetting, settingType string) bool {
	for _, setting := range communitySettings {
		if setting.SettingType == settingType {
			return setting.IsEnabled
		}
	}
	return false
}

// utility method to fetch profile_meta configurations from Cache
func fetchProfileMetaConfigfromCache(redisClient *redis.Client, apiKey string) (*CommunityConfiguration, bool, error) {

	cacheKey := fmt.Sprintf(cache.ProfileMetaConfigurationsCacheKey, apiKey)

	//Fetch profile_meta configurations from cache
	profileMetaValue, exists, err := cache.Get(redisClient, cacheKey)
	if !exists {
		logging.Error(fmt.Sprintf("profile_meta configurations not found in cache for api-key: %s", apiKey))
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
func setProfileMetaConfigInCache(redisClient *redis.Client, apiKey string, profileMetaConfigurations *CommunityConfiguration) error {

	cacheKey := fmt.Sprintf(cache.ProfileMetaConfigurationsCacheKey, apiKey)

	//Save profile_meta configurations in cache
	parsedProfileMeta, err := json.Marshal(profileMetaConfigurations)
	if err != nil {
		return err
	}

	err = cache.Set(redisClient, cacheKey, parsedProfileMeta, cache.ProfileMetaConfigurationsCacheTTL*time.Hour)
	if err != nil {
		logging.Error(fmt.Sprintf("Error while Saving profile_meta configurations in cache for api-key: %s", apiKey))
		return err
	}
	logging.Info(fmt.Sprintf("Saved profile_meta configurations in cache for api-key: %s", apiKey))
	return nil
}
