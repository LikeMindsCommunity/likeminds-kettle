package utils

import (
	"encoding/json"
	"regexp"
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
