package frontendLogger

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/environment"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/logging"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

type AzurePalyloadEntry struct {
	JsonPayload map[string]interface{} `json:"jsonPayload"`
	Timestamp   time.Time              `json:"timestamp"`
}

var (
	appInsightsInitialized bool
	appInsightsClient      appinsights.TelemetryClient
)

func initializeApplicationInsights() {
	instrumentationKey := environment.GoDotEnvVariable("AZURE_APPINSIGHTS_INSTRUMENTATION_KEY")

	if len(instrumentationKey) == 0 {
		logging.Error("Invalid Application Insights Instrumentation Key, Cannot start Application Insights Logger")
		return
	}

	appInsightsClient = appinsights.NewTelemetryClient(instrumentationKey)
	appInsightsInitialized = true
}

func GetAppInsightsClient() appinsights.TelemetryClient {
	return appInsightsClient
}

func logToAppInsights(entries []AzurePalyloadEntry) {
	client := GetAppInsightsClient()

	for _, entry := range entries {
		payload := entry.JsonPayload

		// Extract request and response info
		var method, absoluteURI string
		var httpCode int

		request := appinsights.NewRequestTelemetry(
			method,
			absoluteURI,
			time.Duration(0),
			fmt.Sprint(httpCode),
		)

		// Add meta fields to request.Properties
		for k, v := range payload {
			switch val := v.(type) {
			case time.Time:
				request.Properties[k] = val.Format(time.RFC3339)
			case map[string]interface{}:
				nestedJSON, _ := json.Marshal(val)
				request.Properties[k] = string(nestedJSON)
			default:
				request.Properties[k] = fmt.Sprint(val)
			}
		}

		// Add the IST timestamp to request.Properties
		istLocation, _ := time.LoadLocation("Asia/Kolkata")
		currentTimeIST := time.Now().In(istLocation)
		formattedTimeIST := currentTimeIST.Format("2006-01-02 15:04:05")
		request.Properties["timestamp_IST"] = formattedTimeIST

		// Set success based on status code
		request.Success = httpCode >= 200 && httpCode < 400

		client.Track(request)
	}
}
