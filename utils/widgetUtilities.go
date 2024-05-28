package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/logging"
)

// Response Structure for Custom Widget
type WidgetResponse struct {
	ID               string                 `json:"_id"`
	ParentEntityID   string                 `json:"parent_entity_id"`
	ParentEntityType string                 `json:"parent_entity_type"`
	MetaData         map[string]interface{} `json:"metadata"`
	LMMeta           map[string]interface{} `json:"_lm_meta"`
	CreatedAt        int                    `json:"created_at"`
	UpdatedAt        int                    `json:"updated_at"`
}

type WidgetsResponse struct {
	Success      bool             `json:"success"`
	ErrorMessage string           `json:"error_message"`
	Widgets      []WidgetResponse `json:"widgets"`
}

func fetchWidgetsFromCache(redisClient *redis.Client, communityId int, widgetIds []string) ([]WidgetResponse, []string, error) {

	widgets := []WidgetResponse{}
	remainingWidgetIds := []string{}

	// cache keys for widgets meta
	cachKeys := []string{}
	for _, widgetId := range widgetIds {
		cachKeys = append(cachKeys, fmt.Sprintf(cache.WidgetMetaCacheKey, communityId, widgetId))
	}

	// Fetch keys from cache
	values, err := cache.GetFromMultipleKeys(redisClient, cachKeys...)
	if err != nil {
		return nil, nil, err
	}

	// unmarshall values from cache
	for i, value := range values {
		if value != nil {
			var widget WidgetResponse
			err := json.Unmarshal([]byte(value.(string)), &widget)
			if err != nil {
				remainingWidgetIds = append(remainingWidgetIds, widgetIds[i])
				continue
			}
			widgets = append(widgets, widget)
		} else {
			remainingWidgetIds = append(remainingWidgetIds, widgetIds[i])
		}
	}

	return widgets, remainingWidgetIds, nil
}

func fetchWidgetsFromSwarmService(headers map[string]interface{}, widgetIds []string) ([]WidgetResponse, error) {

	params := map[string]string{
		ParamWidgetIds: ParseStringArrayToString(widgetIds),
	}

	//Send Request
	respBytes, statusCode, err := GetRequestResponseWithoutContext(SwarmService, WidgetEndPoint, GETRequest, headers, params, nil)
	if err != nil || statusCode != http.StatusOK {
		return nil, nil
	}

	//Parse response
	var wr WidgetsResponse
	err = json.Unmarshal(respBytes, &wr)
	if err != nil {
		return nil, err
	}

	if !wr.Success {
		return nil, fmt.Errorf("error fetching widgets meta: %s", wr.ErrorMessage)
	}

	return wr.Widgets, nil
}

func saveWidgetsToCache(redisClient *redis.Client, communityId int, widgets []WidgetResponse) error {

	for _, widget := range widgets {

		parsedWidget, err := json.Marshal(widget)
		if err != nil {
			logging.Error(fmt.Sprintf("error marshalling widget meta: %s", widget.ID))
			continue
		}

		// set widget meta to cache
		cacheKey := fmt.Sprintf(cache.WidgetMetaCacheKey, communityId, widget.ID)
		err = cache.Set(redisClient, cacheKey, parsedWidget, cache.WidgetMetaCacheTTL*time.Hour)
		if err != nil {
			logging.Error(fmt.Sprintf("error setting widget meta to cache: %s", widget.ID))
			return err
		}
		logging.Info(fmt.Sprintf("Widget Meta saved to cache for key: %s", cacheKey))
	}

	return nil
}

// Exposed utility method to fetch widgets from widgetIds from cache if present else from api
func fetchWidgetMetaMapFromWidgetIds(redisClient *redis.Client, headers map[string]interface{}, widgetIds []string,
) (map[string]WidgetResponse, error) {

	// Fetch communityId from ApiKey
	communityId, err := FetchCommunityIdFromApiKey(redisClient, headers[HeadersApiKey].(string))
	if err != nil {
		return nil, err
	}

	widgetsResponse := map[string]WidgetResponse{}

	if len(widgetIds) == 0 {
		return widgetsResponse, nil
	}

	if redisClient != nil {
		// fetch widgets from cache
		widgets, remainingWidgetIds, err := fetchWidgetsFromCache(redisClient, communityId, widgetIds)
		if err != nil {
			return nil, err
		}

		// Add fetched widgets to widgetsResponse
		for _, widget := range widgets {
			widgetsResponse[widget.ID] = widget
		}

		widgetIds = remainingWidgetIds
	}

	if len(widgetIds) > 0 {

		// fetch remaining widgets from api
		widgets, err := fetchWidgetsFromSwarmService(headers, widgetIds)
		if err != nil {
			return nil, err
		}

		// Add fetched widgets to widgetsResponse
		for _, widget := range widgets {
			widgetsResponse[widget.ID] = widget
		}

		// set widgets to cache
		if redisClient != nil {
			go saveWidgetsToCache(redisClient, communityId, widgets)
		}
	}

	return widgetsResponse, nil
}
