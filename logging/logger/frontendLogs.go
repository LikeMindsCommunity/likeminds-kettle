package logger

import (
	"context"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/environment"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type jsonPayload struct {
	Headers       map[string]interface{} `json:"headers"`
	DeviceDetails map[string]interface{} `json:"device_details"`
	StackTrace    map[string]interface{} `json:"stack_trace"`
	SdkMeta       map[string]interface{} `json:"sdk_meta"`
}

type frontendLog struct {
	Timestamp     int64                  `json:"timestamp" binding:"required"`
	DeviceDetails map[string]interface{} `json:"device_details" binding:"required"`
	StackTrace    map[string]interface{} `json:"stack_trace" binding:"required"`
	SdkMeta       map[string]interface{} `json:"sdk_meta"`
	Severity      string                 `json:"severity"`
}

type frontendLogsRequest struct {
	Logs []frontendLog `json:"logs" binding:"required,dive"`
}

func parseFrontendLogsRequest(c *gin.Context) (*frontendLogsRequest, error) {

	var frontendLogs frontendLogsRequest
	if err := c.ShouldBindJSON(&frontendLogs); err != nil {
		return &frontendLogs, err
	}

	return &frontendLogs, nil
}

func createHeadersForLogging(c *gin.Context) map[string]interface{} {

	headers := map[string]interface{}{
		"member_id":     user.GetRequestingUserId(c),
		"api_key":       c.GetHeader(utils.HeadersApiKey),
		"platform_code": c.GetHeader(utils.HeadersPlatformCode),
		"version_code":  c.GetHeader(utils.HeadersVersionCode),
		"sdk_source":    c.GetHeader(utils.HeadersSdkSource),
	}

	return headers
}

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

func getFrontendLogger() (*logging.Client, *logging.Logger, error) {

	projectId := environment.GoDotEnvVariable("GCP_LOGGING_PROJECT_ID")

	ctx := context.Background()

	logClient, err := logging.NewClient(ctx, projectId)
	if err != nil {
		return nil, nil, err
	}

	logger := logClient.Logger(FrontendLogger)

	return logClient, logger, nil
}

func pushLogsInternal(headers map[string]interface{}, logger *logging.Logger, logs []frontendLog) error {

	logPayload := jsonPayload{
		Headers: headers,
	}

	for _, log := range logs {

		// Update jsonPayload
		logPayload.DeviceDetails = log.DeviceDetails
		logPayload.StackTrace = log.StackTrace
		logPayload.SdkMeta = log.SdkMeta

		// Create log entry
		logEntry := logging.Entry{
			Timestamp: time.Unix(log.Timestamp, 0),
			Payload:   logPayload,
			Severity:  getSeverityFromLogLevel(log.Severity),
		}

		// Send log entry
		logger.Log(logEntry)
	}

	return nil
}

func PushFrontendLogs(c *gin.Context) {

	// Parse request
	flr, err := parseFrontendLogsRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Get logger for pushing frontend logs
	logClient, logger, err := getFrontendLogger()
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Close client when done.
	defer logClient.Close()

	// Push logs
	err = pushLogsInternal(createHeadersForLogging(c), logger, flr.Logs)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	utils.GenerateResponse(c, nil, false)
}
