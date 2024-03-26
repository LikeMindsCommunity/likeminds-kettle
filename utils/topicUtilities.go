package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/logging"
)

// TopicMeta | schema for topic meta
type TopicMeta struct {
	ID           string   `json:"_id"`
	Name         string   `json:"name"`
	IsEnabled    bool     `json:"is_enabled"`
	Priority     float32  `json:"priority"`
	IsSearchable bool     `json:"is_searchable"`
	ParentId     string   `json:"parent_id"`
	ParentName   string   `json:"parent_name"`
	AllParentIds []string `json:"all_parent_ids"`
	Level        int      `json:"level"`
	WidgetId     string   `json:"widget_id"`
}

type FetchTopicsResponse struct {
	Success      bool                      `json:"success"`
	ErrorMessage string                    `json:"error_message"`
	Topics       []TopicMeta               `json:"topics"`
	Widgets      map[string]WidgetResponse `json:"widgets"`
}

func fetchTopicsFromCache(redisClient *redis.Client, communityId int, topicIds []string) ([]TopicMeta, []string, error) {

	topicsMeta := []TopicMeta{}
	remainingTopicIds := []string{}

	// cache keys for topics meta
	cachKeys := []string{}
	for _, topicId := range topicIds {
		cachKeys = append(cachKeys, fmt.Sprintf(cache.TopicMetaCacheKey, communityId, topicId))
	}

	// Fetch keys from cache
	values, err := cache.GetFromMultipleKeys(redisClient, cachKeys...)
	if err != nil {
		return nil, nil, err
	}

	// unmarshall values from cache
	for i, value := range values {
		if value != nil {
			var topicMeta TopicMeta
			err := json.Unmarshal([]byte(value.(string)), &topicMeta)
			if err != nil {
				logging.Error(fmt.Sprint("Error unmarshalling topic meta from cache", err))
				continue
			}

			topicsMeta = append(topicsMeta, topicMeta)
		} else {
			remainingTopicIds = append(remainingTopicIds, topicIds[i])
		}
	}

	return topicsMeta, remainingTopicIds, nil
}

func fetchTopicsFromSwarmService(headers map[string]interface{}, topicIds []string,
) ([]TopicMeta, map[string]WidgetResponse, error) {

	// Fetch topics meta from swarm service
	params := map[string]string{
		ParamParentIds: ParseStringArrayToString(topicIds),
		ParamPageSize:  "1",
	}

	// Send Request
	respBytes, _, err := GetRequestResponseWithoutContext(SwarmService, FetchTopicsEndpoint, GETRequest, headers, params, nil)
	if err != nil {
		return nil, nil, err
	}

	// Parse response
	var tr FetchTopicsResponse
	err = json.Unmarshal(respBytes, &tr)
	if err != nil {
		return nil, nil, err
	}

	if !tr.Success {
		return nil, nil, fmt.Errorf("error fetching topics meta: %s", tr.ErrorMessage)
	}

	return tr.Topics, tr.Widgets, nil
}

func saveTopicsInCache(redisClient *redis.Client, communityId int, topicsMeta []TopicMeta) {

	for _, topicMeta := range topicsMeta {

		parsedData, err := json.Marshal(topicMeta)
		if err != nil {
			logging.Error(fmt.Sprint("error marshalling topic data", err))
			continue
		}

		// save to cache
		cacheKey := fmt.Sprintf(cache.TopicMetaCacheKey, communityId, topicMeta.ID)
		err = cache.Set(redisClient, cacheKey, parsedData, cache.TopicMetaCacheTTL*time.Hour)
		if err != nil {
			logging.Error(fmt.Sprint("error saving topic meta to cache", err))
		}
		logging.Info(fmt.Sprintf("Topic meta saved to cache with key: %s", cacheKey))
	}
}

// External utility method to fetch topics meta map from topic ids from cache if present else from API
func FetchTopicsMetaFromTopicsIds(redisClient *redis.Client, headers map[string]interface{}, topicIds []string) (map[string]TopicMeta, error) {

	// Fetch communityId from ApiKey
	communityId, err := FetchCommunityIdFromApiKey(redisClient, headers[HeadersApiKey].(string))
	if err != nil {
		return nil, err
	}

	topicsMeta := map[string]TopicMeta{}

	if len(topicIds) == 0 {
		return topicsMeta, nil
	}

	if redisClient != nil {

		// Fetch topics meta from cache
		cachedTopicsMeta, remainingTopicIds, err := fetchTopicsFromCache(redisClient, communityId, topicIds)
		if err != nil {
			return nil, err
		}

		// convert topics meta to map
		for _, topicMeta := range cachedTopicsMeta {
			topicsMeta[topicMeta.ID] = topicMeta
		}

		topicIds = remainingTopicIds
	}

	// Fetch topics meta from swarm service
	if len(topicIds) > 0 {

		fetchedTopicsMeta, fetchedWidgetsMeta, err := fetchTopicsFromSwarmService(headers, topicIds)
		if err != nil {
			return nil, err
		}

		if redisClient != nil {
			// save fetched topics meta to cache
			go saveTopicsInCache(redisClient, communityId, fetchedTopicsMeta)

			// save fetched widgets meta to cache
			go func() {
				widgetsResponse := []WidgetResponse{}
				for _, widgetResponse := range fetchedWidgetsMeta {
					widgetsResponse = append(widgetsResponse, widgetResponse)
				}

				saveWidgetsToCache(redisClient, communityId, widgetsResponse)
			}()
		}

		// convert topics meta to map
		for _, topicMeta := range fetchedTopicsMeta {
			topicsMeta[topicMeta.ID] = topicMeta
		}
	}

	return topicsMeta, nil
}
