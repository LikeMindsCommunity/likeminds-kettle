package utils

import (
	"encoding/json"
	"net/http"
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
	Success bool             `json:"success"`
	Widgets []WidgetResponse `json:"widgets"`
}

// Exposed utility method to fetch widgets from widgetIds
func GetWidgetsFromWidgetIds(headers map[string]interface{}, widgetIds []string) ([]WidgetResponse, error) {

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

	return wr.Widgets, nil
}
