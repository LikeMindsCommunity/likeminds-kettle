package frontendLogger

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/environment"
	log "github.com/nateshr/likeminds-authentication/internal/logging"
)

var (
	logClient *logging.Client
)

// PayloadEntry is the struct for the payload to be sent to the logger
type PayloadEntry struct {
	JsonPayload map[string]interface{} `json:"jsonPayload"`
	Severity    string                 `json:"severity"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Internal method to get logger with logger name
func getLogger(logId string) (*logging.Logger, error) {

	// Initialize logClient
	projectId := environment.GoDotEnvVariable("GCP_LOGGING_PROJECT_ID")
	ctx := context.Background()

	if logClient == nil {
		var err error
		logClient, err = logging.NewClient(ctx, projectId)
		if err != nil {
			log.Error(fmt.Sprintf("Error creating new LogClient => %s", err.Error()))
			return nil, err
		}
	}

	// Sets the name of the log to write to.
	logger := logClient.Logger(logId)

	return logger, nil
}

// Exposed method get Frontend Logger
func GetFrontendLogger() (*logging.Logger, error) {

	return getLogger(FrontendLoggerId)
}

// Internal method to get severity from log level
func getSeverityFromLogLevel(logLevel string) logging.Severity {

	switch logLevel {
	case "debug":
		return logging.Debug
	case "info":
		return logging.Info
	case "notice":
		return logging.Notice
	case "warning":
		return logging.Warning
	case "error":
		return logging.Error
	case "critical":
		return logging.Critical
	case "alert":
		return logging.Alert
	case "emergency":
		return logging.Emergency
	default:
		return logging.Default
	}
}

// Exposed method to push logs using logger
func PushLogEntries(entries []PayloadEntry, logger *logging.Logger) {

	for _, entry := range entries {

		// Create log entry
		logEntry := logging.Entry{
			Timestamp: entry.Timestamp,
			Payload:   entry.JsonPayload,
			Severity:  getSeverityFromLogLevel(entry.Severity),
		}

		// Send log entry
		logger.Log(logEntry)
	}
}

// Exposed method to get all the headers for logging
func GetHeadersForLogging(c *gin.Context) map[string]interface{} {

	// map to store headers for logging
	headers := make(map[string]interface{})

	for key, value := range c.Request.Header {
		headers[key] = value[0]
	}

	return headers
}
