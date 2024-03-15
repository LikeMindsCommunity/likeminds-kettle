package utils

import (
	"encoding/json"
	"regexp"

	"github.com/gin-gonic/gin"
)

// Exposed utility method to parse response for widget_ids using regex
func GetProfileWidgetIdsFromDataResponse(dataResponse map[string]interface{}) []string {

	var widgetIds []string

	// marshall the dataResponse to string
	inputText, _ := json.Marshal(dataResponse)

	re := regexp.MustCompile(`['"]sdk_client_info['"]\s*:\s*{[^{}]*['"]widget_id['"]\s*:\s*['"]([^'"]+)`)
	matches := re.FindAllStringSubmatch(string(inputText), -1)

	for _, match := range matches {
		if len(match) > 1 {
			widgetIds = append(widgetIds, match[1])
		}
	}

	return widgetIds
}

func ParseAndFetchProfileWidgets(c *gin.Context, userId string, dataResponse map[string]interface{}) map[string]interface{} {

	if userId == "" {
		return dataResponse
	}

	// Fetch Profile meta configurations and check if widget are enabled
	profileWidgetsEnabled, _ := IsProfileWidgetsEnabled(c, userId)
	if profileWidgetsEnabled {

		if dataResponse["widgets"] == nil {
			dataResponse["widgets"] = map[string]interface{}{}
		}

		// parse and fetch widget ids from data response
		widgetIds := GetProfileWidgetIdsFromDataResponse(dataResponse)

		if len(widgetIds) > 0 {

			// fetch widgets meta for the widget ids
			widgetsMap, _ := fetchWidgetMetaMapFromWidgetIds(GetRedisClientFromContext(c), CreateHeaders(c, userId), widgetIds)

			for widgetId, widgetMeta := range widgetsMap {
				dataResponse["widgets"].(map[string]interface{})[widgetId] = widgetMeta
			}
		}

	}

	return dataResponse
}
