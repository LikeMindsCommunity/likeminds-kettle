package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/logging"
)

type UserTopics map[string][]string

type UserTopicsResponse struct {
	Success      bool                      `json:"success"`
	ErrorMessage string                    `json:"error_message"`
	UserTopics   UserTopics                `json:"user_topics"`
	Topics       map[string]TopicMeta      `json:"topics"`
	Widgets      map[string]WidgetResponse `json:"widgets"`
}

func fetchUserTopicsFromCache(redisClient *redis.Client, userUniqueIds []string) (UserTopics, []string, error) {

	userTopics := UserTopics{}
	remainingUserUniqueIds := []string{}

	// cache keys for user topics
	userTopicsCacheKeys := []string{}
	for _, userUniqueId := range userUniqueIds {
		userTopicsCacheKeys = append(userTopicsCacheKeys, fmt.Sprintf(cache.UserTopicsCacheKey, userUniqueId))
	}

	// Fetch keys from cache
	values, err := cache.MGet(redisClient, userTopicsCacheKeys...)
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

func saveUserTopicsToCache(redisClient *redis.Client, userTopics UserTopics) error {

	for userUniqueId, topics := range userTopics {
		parsedTopics, err := json.Marshal(topics)
		if err != nil {
			logging.Error(fmt.Sprintf("Error marshalling user topics: %s", err))
			continue
		}

		cacheKey := fmt.Sprintf(cache.UserTopicsCacheKey, userUniqueId)
		err = cache.Set(redisClient, cacheKey, parsedTopics, cache.UserTopicsCacheTTL*time.Hour) // TODO: Move to constants
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

	userTopicsMap := UserTopics{}

	if len(userUniqueIds) == 0 {
		return userTopicsMap, nil
	}

	if redisClient != nil {
		userTopics, remainingUserUniqueIds, err := fetchUserTopicsFromCache(redisClient, userUniqueIds)
		if err != nil {
			return userTopics, err
		}

		userTopicsMap = userTopics

		userUniqueIds = remainingUserUniqueIds
	}

	// fetch user topics from API
	if len(userUniqueIds) > 0 {

		// fetch user topics from API
		userTopicsFromAPI, topicsFromAPI, widgetsFromAPI, err := fetchUserTopicsFromSwarmService(headers, userUniqueIds)
		if err != nil {
			return nil, err
		}

		if redisClient != nil {

			// save user topics to cache
			go saveUserTopicsToCache(redisClient, userTopicsFromAPI)

			// save topics meta to cache
			go func() {
				topicsData := []TopicMeta{}
				for _, topic := range topicsFromAPI {
					topicsData = append(topicsData, topic)
				}

				// save topics to cache
				saveTopicsInCache(redisClient, topicsData)
			}()

			// save widgets meta to cache
			go func() {
				widgetsData := []WidgetResponse{}
				for _, widget := range widgetsFromAPI {
					widgetsData = append(widgetsData, widget)
				}

				// save widgets to cache
				saveWidgetsToCache(redisClient, widgetsData)
			}()

		}

		// merge user topics from cache and API
		for userUniqueId, topics := range userTopicsFromAPI {
			userTopicsMap[userUniqueId] = topics
		}
	}

	return userTopicsMap, nil
}
