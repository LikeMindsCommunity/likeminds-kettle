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
			for _, report := range reports.([]map[string]interface{}) {
				if typeValue, ok := report["type"]; ok {
					if typeValue == feed.POST_REPORT_TYPE {
						post_data := feed.GetPostInternal(c, userId, report["entity_id"].(string))
						if post_data != nil {
							report["entity_data"] = post_data
						}
					}

					if typeValue == feed.COMMENT_REPORT_TYPE || typeValue == feed.REPLY_REPORT_TYPE {
						comment_data := feed.FetchCommentByIdInternal(c, userId, report["entity_id"].(string))
						if comment_data != nil {
							report["entity_data"] = comment_data
						}
					}
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
