package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/logging"
)

type UserTopics map[string][]string

type UserTopicsResponse struct {
	Success      bool                      `json:"success"`
	ErrorMessage string                    `json:"error_message"`
	UserTopics   UserTopics                `json:"user_topics"`
	Topics       map[string]TopicMeta      `json:"topics"`
	Widgets      map[string]WidgetResponse `json:"widgets"`
}

func fetchUserTopicsFromCache(redisClient *redis.Client, communityId int, userUniqueIds []string) (UserTopics, []string, error) {

	userTopics := UserTopics{}
	remainingUserUniqueIds := []string{}

	// cache keys for user topics
	userTopicsCacheKeys := []string{}
	for _, userUniqueId := range userUniqueIds {
		userTopicsCacheKeys = append(userTopicsCacheKeys, fmt.Sprintf(cache.UserTopicsCacheKey, communityId, userUniqueId))
	}

	// Fetch keys from cache
	values, err := cache.GetFromMultipleKeys(redisClient, userTopicsCacheKeys...)
	if err != nil {
		return nil, nil, err
	}

	// unmarshall values from cache
	for i, value := range values {
		if value != nil {
			var topics []string
			err := json.Unmarshal([]byte(value.(string)), &topics)
			if err != nil {
				remainingUserUniqueIds = append(remainingUserUniqueIds, userUniqueIds[i])
				continue
			}
			userTopics[userUniqueIds[i]] = topics
		} else {
			remainingUserUniqueIds = append(remainingUserUniqueIds, userUniqueIds[i])
		}
	}

	return userTopics, remainingUserUniqueIds, nil
}

func fetchUserTopicsFromSwarmService(headers map[string]interface{}, userUniqueIds []string,
) (UserTopics, map[string]TopicMeta, map[string]WidgetResponse, error) {

	params := map[string]string{
		ParamUUIDs: ParseStringArrayToString(userUniqueIds),
	}

	// Send Request
	respBytes, _, err := GetRequestResponseWithoutContext(SwarmService, FetchUserTopicsEndpoint, GETRequest, headers, params, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	// Parse response
	var utr UserTopicsResponse
	err = json.Unmarshal(respBytes, &utr)
	if err != nil {
		return nil, nil, nil, err
	}

	if !utr.Success {
		return nil, nil, nil, fmt.Errorf(utr.ErrorMessage)
	}

	return utr.UserTopics, utr.Topics, utr.Widgets, nil
}

func saveUserTopicsToCache(redisClient *redis.Client, communityId int, userTopics UserTopics) error {

	for userUniqueId, topics := range userTopics {
		parsedTopics, err := json.Marshal(topics)
		if err != nil {
			logging.Error(fmt.Sprintf("Error marshalling user topics: %s", err))
			continue
		}

		cacheKey := fmt.Sprintf(cache.UserTopicsCacheKey, communityId, userUniqueId)
		err = cache.Set(redisClient, cacheKey, parsedTopics, cache.UserTopicsCacheTTL*time.Hour)
		if err != nil {
			logging.Error(fmt.Sprintf("Error saving user topics to cache: %s", err))
			return err
		}
		logging.Info(fmt.Sprintf("User topics saved to cache with key: %s", cacheKey))
	}
	return nil
}

// Exposed utility method to fetch user topics from userUniqueIds from cache if present else from API
func FetchUserTopicsForUserUniqueIds(redisClient *redis.Client, headers map[string]interface{}, userUniqueIds []string,
) (UserTopics, error) {

	// Fetch communityId from ApiKey
	communityId, err := FetchCommunityIdFromApiKey(redisClient, headers[HeadersApiKey].(string))
	if err != nil {
		return nil, err
	}

	userTopicsMap := UserTopics{}

	if len(userUniqueIds) == 0 {
		return userTopicsMap, nil
	}

	if redisClient != nil {
		userTopics, remainingUserUniqueIds, err := fetchUserTopicsFromCache(redisClient, communityId, userUniqueIds)
		if err != nil {
			return userTopics, err
		}

		userTopicsMap = userTopics

		userUniqueIds = remainingUserUniqueIds
	}

	// fetch user topics from API
	if len(userUniqueIds) > 0 {

		// fetch user topics from API
		userTopicsFromAPI, err := fetchUserTopicsFromApiAndSaveInCache(headers, userUniqueIds, redisClient, communityId)
		if err != nil {
			return nil, err
		}

		// merge user topics from cache and API
		for userUniqueId, topics := range userTopicsFromAPI {
			userTopicsMap[userUniqueId] = topics
		}
	}

	return userTopicsMap, nil
}

func fetchUserTopicsFromApiAndSaveInCache(headers map[string]interface{}, userUniqueIds []string, redisClient *redis.Client, communityId int,
) (UserTopics, error) {

	userTopics, topicsFromAPI, widgetsFromAPI, err := fetchUserTopicsFromSwarmService(headers, userUniqueIds)
	if err != nil {
		return nil, err
	}

	if redisClient != nil {

		// save user topics to cache
		SafeGo(func() { saveUserTopicsToCache(redisClient, communityId, userTopics) })

		// save topics meta to cache
		SafeGo(func() {
			topicsData := []TopicMeta{}
			for _, topic := range topicsFromAPI {
				topicsData = append(topicsData, topic)
			}

			// save topics to cache
			saveTopicsInCache(redisClient, communityId, topicsData)
		})

		// save widgets meta to cache
		SafeGo(func() {
			widgetsData := []WidgetResponse{}
			for _, widget := range widgetsFromAPI {
				widgetsData = append(widgetsData, widget)
			}

			// save widgets to cache
			saveWidgetsToCache(redisClient, communityId, widgetsData)
		})

	}

	return userTopics, nil
}

// External method to fetch user topics (If enabled) and its related data for userUniqueIds and update in dataResponse
func FetchAndUpdateUserTopicsDataForResponse(redisClient *redis.Client, headers map[string]interface{}, dataResponse map[string]interface{}, userUniqueIds []string,
) map[string]interface{} {

	defer Timer("FetchAndUpdateUserTopicsDataForResponse")()

	userTopics, topicsMeta, widgetsMeta := fetchUserTopicsAndWidgetsMeta(redisClient, headers, userUniqueIds)

	// update userTopics in dataResponse
	if dataResponse[constants.ResponseKeyUserTopics] == nil {
		dataResponse[constants.ResponseKeyUserTopics] = map[string]interface{}{}
	}
	for key, value := range userTopics {
		dataResponse[constants.ResponseKeyUserTopics].(map[string]interface{})[key] = value
	}

	// Update topic meta in dataResponse
	if dataResponse[constants.ResponseKeyTopics] == nil {
		dataResponse[constants.ResponseKeyTopics] = map[string]interface{}{}
	}
	for key, value := range topicsMeta {
		dataResponse[constants.ResponseKeyTopics].(map[string]interface{})[key] = value
	}

	// Update widget meta in dataResponse
	if dataResponse[constants.ResponseKeyWidgets] == nil {
		dataResponse[constants.ResponseKeyWidgets] = map[string]interface{}{}
	}
	for key, value := range widgetsMeta {
		dataResponse[constants.ResponseKeyWidgets].(map[string]interface{})[key] = value
	}

	return dataResponse
}

// method to fetch user topics and widgets meta
func fetchUserTopicsAndWidgetsMeta(redisClient *redis.Client, headers map[string]interface{},
	userUniqueIds []string) (UserTopics, map[string]TopicMeta, map[string]WidgetResponse) {

	userTopics, topicsMeta, widgetsMeta := UserTopics{}, map[string]TopicMeta{}, map[string]WidgetResponse{}

	// if user Topics connection is enabled, fetch related data
	if UserTopicsConnectionEnabled(redisClient, headers) {
		err := error(nil)

		// fetch user topics data for user_unique_ids
		userTopics, err = FetchUserTopicsForUserUniqueIds(redisClient, headers, userUniqueIds)
		if err != nil {
			logging.Error(fmt.Sprint("error fetching user topics for user_unique_ids", err))
		}

		// Fetch topics meta for user topics
		topicsIds := []string{}
		for _, userTopics := range userTopics {
			topicsIds = append(topicsIds, userTopics...)
		}
		topicsMeta, err = FetchTopicsMetaFromTopicsIds(redisClient, headers, topicsIds)
		if err != nil {
			logging.Error(fmt.Sprint("error fetching topics meta for topics ids", err))
		}

		// Fetch widget meta for user topics
		widgetIds := []string{}
		for _, userTopic := range topicsMeta {
			if userTopic.WidgetId != "" {
				widgetIds = append(widgetIds, userTopic.WidgetId)
			}
		}
		widgetsMeta, err = fetchWidgetMetaMapFromWidgetIds(redisClient, headers, widgetIds)
		if err != nil {
			logging.Error(fmt.Sprint("error fetching widget meta for widget ids", err))
		}
	}

	return userTopics, topicsMeta, widgetsMeta
}
