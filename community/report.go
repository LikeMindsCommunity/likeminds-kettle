package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PushReportRequest struct {
	ReportedMemberID int    `json:"reported_member_id"`
	ConversationID   int    `json:"conversation_id"`
	ChatroomID       int    `json:"collabcard_id"`
	EntityId         string `json:"entity_id"`
	EntityCreatorId  string `json:"entity_creator_id"`
	EntityType       int    `json:"entity_type"`
	Link             string `json:"link"`
	TagId            int    `json:"tag_id"`
	Reason           string `json:"reason"`
}

type CloseReportRequest struct {
	ReportID int `json:"report_id" binding:"required"`
}

// GetReport is used to get community reports
func GetReport(c *gin.Context) {
	Report(c, utils.GETMethod)
}

// PushReport is used to push a report in a community
func PushReport(c *gin.Context) {
	Report(c, utils.POSTMethod)
}

// DeleteReport is used to delete a report in community
func CloseReport(c *gin.Context) {
	Report(c, utils.DELETEMethod)
}

// Report method handles community reports
func Report(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:

		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		//Send Request
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchReportsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
		if respBytes == nil {
			return
		}

		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return
		}

		//If flow succeeds
		dataResponse := apiCR.Response
		if reports, ok := dataResponse["reports"]; ok {
			for _, report := range reports.([]interface{}) {
				report := fetchReportEntityData(c, report, userId)
				if report == nil {
					return
				}
			}
		}

		//Generate response
		utils.GenerateResponse(c, dataResponse)

	case utils.POSTMethod:

		pushReportRequest, err := parsePushReportRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, PushReportEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, pushReportRequest)

	case utils.DELETEMethod:

		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		//Body to be sent in the close report api internally
		closeReportRequest, err := parseCloseReportRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, CloseReportsEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, closeReportRequest)

	}
}

func parsePushReportRequest(c *gin.Context) (*PushReportRequest, error) {
	//POST body params
	var prr PushReportRequest

	if err := c.ShouldBindJSON(&prr); err != nil {
		return nil, err
	}

	return &prr, nil
}

func parseCloseReportRequest(c *gin.Context) (*CloseReportRequest, error) {
	//POST body params
	var crr CloseReportRequest

	if err := c.ShouldBindJSON(&crr); err != nil {
		return nil, err
	}

	return &crr, nil
}

func fetchReportEntityData(c *gin.Context, report interface{}, userId string) interface{} {
	typeValue, ok := report.(map[string]interface{})["type"]
	if ok {

		headers := utils.CreateHeaders(c, userId)
		params := make(map[string]string)

		if int(typeValue.(float64)) == feed.POST_REPORT_TYPE {

			entity_id := report.(map[string]interface{})["entity_id"].(string)

			// Fetch post data without context
			post_data, err := feed.GetPostWithoutContext(headers, params, entity_id, true)

			// If found post data, then populate entity_data, else set it to nil
			if err != nil {
				report.(map[string]interface{})["entity_data"] = nil
			} else {
				report.(map[string]interface{})["entity_data"] = post_data
			}
		}

		if int(typeValue.(float64)) == feed.COMMENT_REPORT_TYPE || int(typeValue.(float64)) == feed.REPLY_REPORT_TYPE {

			entity_id := report.(map[string]interface{})["entity_id"].(string)

			// Fetch comment data without context
			comment_data, err := feed.GetCommentWithoutContext(headers, params, entity_id, true)

			// If found comment data, then populate entity_data, else set it to nil
			if err != nil {
				report.(map[string]interface{})["entity_data"] = nil
				return nil
			} else {
				report.(map[string]interface{})["entity_data"] = comment_data

			}

		}
	}

	return report
}
