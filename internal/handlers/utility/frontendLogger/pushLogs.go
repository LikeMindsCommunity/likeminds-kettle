package frontendLogger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/environment"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type logRequest struct {
	Timestamp  int64                  `json:"timestamp"`
	DeviceMeta map[string]interface{} `json:"device_meta" binding:"required"`
	StackTrace map[string]interface{} `json:"stack_trace" binding:"required"`
	SdkMeta    map[string]interface{} `json:"sdk_meta"`
	Severity   string                 `json:"severity"`
}

type logsRequest struct {
	Logs []logRequest `json:"logs" binding:"required,dive"`
}

func parseLogsRequest(c *gin.Context) (*logsRequest, error) {

	var lr logsRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		return &lr, err
	}

	return &lr, nil
}

func PushLogs(c *gin.Context) {
	// Parse request
	flr, err := parseLogsRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Check if platform code is valid
	platform_code := c.GetHeader(utils.HeadersPlatformCode)
	if !utils.CheckIfStringExistsInArray(utils.ValidFrontendPlatformCodes, platform_code) {
		utils.GeneralBadRequestError(c, fmt.Sprintf("Invalid platform code: %s", platform_code))
		return
	}

	headers := GetHeadersForLogging(c)

	logPlatform := environment.GoDotEnvVariable(environment.EnvLogPlatform)

	if logPlatform == "GCP" {
		pushToGCP(c, headers, flr)
	} else {
		// default to cloudwatch
		err = pushToCloudwatch(platform_code, headers, flr)
		if err != nil {
			utils.GeneralAPIError(c, fmt.Sprint("cloudwatch error: ", err))
			return
		}
	}

	// Send response
	utils.GenerateResponse(c, map[string]interface{}{}, false)
}

func pushToCloudwatch(platformCode string, headers map[string]interface{}, flr *logsRequest) error {
	// Create CloudWatchLogs client
	client, err := GetCloudwatchClient()
	if err != nil {
		logging.Error(fmt.Sprint("Error loading cloudwatch client: ", err.Error()))
		return err
	}

	// Define log group and stream names
	logGroupName := LogGroupName
	logStreamName := utils.GeneratePlatformString(platformCode)

	// Create log group if not exists
	if err := CreateLogGroupIfNotExist(client, logGroupName); err != nil {
		logging.Error(fmt.Sprint("Already exists: ", err.Error()))
	}

	// Create log stream if not exists
	if err := CreateLogStreamIfNotExist(client, logGroupName, logStreamName); err != nil {
		logging.Error(fmt.Sprint("Already exists: ", err.Error()))
	}

	entries := createPayloadEntriesCloudwatch(headers, flr.Logs)

	// Log a message to CloudWatch
	err = LogToCloudWatch(client, logGroupName, logStreamName, entries)
	if err != nil {
		logging.Error(fmt.Sprint("Some error occured while pushing to cloudwatch: ", err.Error()))
		return err
	}
	return nil
}

func pushToGCP(c *gin.Context, headers map[string]interface{}, flr *logsRequest) {
	// Get front logger for pushing frontend logs
	logger, err := GetFrontendLogger()
	if err != nil {
		// we'll be sending 200 until we have implemented another library
		logging.Error(err.Error())
		utils.GenerateResponse(c, map[string]interface{}{}, false)
		return
	}

	// Create payload entries
	entries := createPayloadEntries(headers, flr.Logs)

	// Push log entries
	PushLogEntries(entries, logger)
}

// Get valid timestamp for logging
func getValidTimestamp(logTimestamp int64) (time.Time, error) {

	timestamp := time.Now()

	if logTimestamp > 0 {

		timestampLowerlimit := time.Now().AddDate(0, -1, 0).Unix() * 1000 // 1 month in past
		timestampUpperlimit := time.Now().AddDate(0, 0, 1).Unix() * 1000  // 1 day in future

		// check if timestamp is valid
		if (logTimestamp < timestampLowerlimit) || (logTimestamp > timestampUpperlimit) {
			return timestamp, fmt.Errorf("invalid timestamp: %d", logTimestamp)
		}

		timestamp = time.Unix(0, logTimestamp*int64(time.Millisecond))
	}

	return timestamp, nil
}

// create payload for entries
func createPayloadEntries(headers map[string]interface{}, logs []logRequest) []PayloadEntry {

	var entries []PayloadEntry

	// create payload entry for each log
	for _, lr := range logs {

		payload := map[string]interface{}{
			"device_details": lr.DeviceMeta,
			"stack_trace":    lr.StackTrace,
			"sdk_meta":       lr.SdkMeta,
			"headers":        headers,
		}

		timestamp, err := getValidTimestamp(lr.Timestamp)
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		entry := PayloadEntry{
			JsonPayload: payload,
			Severity:    lr.Severity,
			Timestamp:   timestamp,
		}

		entries = append(entries, entry)
	}

	return entries

}

func createPayloadEntriesCloudwatch(headers map[string]interface{}, logs []logRequest) []CloudwatchPayloadEntry {

	var entries []CloudwatchPayloadEntry

	// create payload entry for each log
	for _, lr := range logs {

		payload := map[string]interface{}{
			"device_details": lr.DeviceMeta,
			"stack_trace":    lr.StackTrace,
			"sdk_meta":       lr.SdkMeta,
			"headers":        headers,
			"severity":       lr.Severity,
		}

		timestamp, err := getValidTimestamp(lr.Timestamp)
		if err != nil {
			logging.Error(err.Error())
			continue
		}

		entry := CloudwatchPayloadEntry{
			JsonPayload: payload,
			Timestamp:   timestamp,
		}

		entries = append(entries, entry)
	}

	return entries
}
