package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/utils"
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

	// Get front logger for pushing frontend logs
	logger, err := logging.GetFrontendLogger()
	if err != nil {
		// we'll be sending 200 until we have implemented another library
		logging.Error(err.Error())
		utils.GenerateResponse(c, map[string]interface{}{}, false)
		return
	}

	// Get headers dumps from request
	headers := logging.GetHeadersForLogging(c)

	// Create payload entries
	entries := createPayloadEntries(headers, flr.Logs)

	// Push log entries
	logging.PushLogEntries(entries, logger)

	// Send response
	utils.GenerateResponse(c, map[string]interface{}{}, false)
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
func createPayloadEntries(headers map[string]interface{}, logs []logRequest) []logging.PayloadEntry {

	var entries []logging.PayloadEntry

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

		entry := logging.PayloadEntry{
			JsonPayload: payload,
			Severity:    lr.Severity,
			Timestamp:   timestamp,
		}

		entries = append(entries, entry)
	}

	return entries

}
