package utils

import (
	"encoding/json"
	"regexp"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/gin-gonic/gin"
)

// Exposed utility method to parse response for widget_ids using regex
func GetWidgetIdsFromDataResponse(dataResponse map[string]interface{}) []string {

	var widgetIds []string

	// marshall the dataResponse to string
	inputText, _ := json.Marshal(dataResponse)

	re := regexp.MustCompile(`['"]widget_id['"]\s*:\s*['"]([^'"]+)`)
	matches := re.FindAllStringSubmatch(string(inputText), -1)

	for _, match := range matches {
		if len(match) > 1 {
			widgetIds = append(widgetIds, match[1])
		}
	}

	return widgetIds
}

func ParseAndFetchWidgets(c *gin.Context, userId string, dataResponse map[string]interface{}) map[string]interface{} {

	if userId == "" {
		return dataResponse
	}

	if dataResponse[constants.ResponseKeyWidgets] == nil {
		dataResponse[constants.ResponseKeyWidgets] = map[string]interface{}{}
	}

	// Parse and fetch widget ids from data response
	widgetIds := GetWidgetIdsFromDataResponse(dataResponse)

	if len(widgetIds) > 0 {

		// fetch widgets meta for the widget ids
		widgetsMap, _ := fetchWidgetMetaMapFromWidgetIds(GetRedisClientFromContext(c), CreateHeaders(c, userId), widgetIds)

		for widgetId, widgetMeta := range widgetsMap {
			dataResponse[constants.ResponseKeyWidgets].(map[string]interface{})[widgetId] = widgetMeta
		}
	}

	return dataResponse
}
