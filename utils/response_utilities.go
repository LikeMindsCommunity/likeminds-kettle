package utils

import (
	"encoding/json"
	"regexp"

	"github.com/gin-gonic/gin"
)

// Exposed utility method to parse response for widget_ids using regex
func GetWidgetIdsFromDataResponse(dataResponse map[string]interface{}) []string {

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
	profileWidgetsEnabled, _ := ProfileWidgetsEnabled(c, userId)
	if profileWidgetsEnabled {

		if dataResponse["widgets"] == nil {
			dataResponse["widgets"] = map[string]interface{}{}
		}

		// parse and fetch widget ids from data response
		widgetIds := GetWidgetIdsFromDataResponse(dataResponse)

		if len(widgetIds) > 0 {

			// fetch widgets from widget ids from swarm service
			widgets, _ := GetWidgetsFromWidgetIds(CreateHeaders(c, userId), widgetIds)

			for _, value := range widgets {
				widgetId := value.ID
				dataResponse["widgets"].(map[string]interface{})[widgetId] = value
			}
		}

	}

	return dataResponse
}
