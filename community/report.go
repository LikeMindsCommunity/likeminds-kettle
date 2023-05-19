package community

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PushReportRequest struct {
	ReportedMemberID int         `json:"reported_member_id"`
	ConversationID   interface{} `json:"conversation_id"`
	CollabcardId     interface{} `json:"collabcard_id"`
	EntityId         string      `json:"entity_id"`
	EntityCreatorId  string      `json:"entity_creator_id"`
	EntityType       int         `json:"entity_type"`
	Link             string      `json:"link"`
	TagId            int         `json:"tag_id"`
	Reason           string      `json:"reason"`
	ChatroomId       interface{} `json:"chatroom_id"`
	UUID             string      `json:"uuid"`
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

		getReportsInternal(c, userId)

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

func getReportsInternal(c *gin.Context, userId string) {

	// Get Bot Id if request from dashboard
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Params to be sent with pagination and filter support in API
	params := map[string]string{
		ParamPage:       c.Query(ParamPage),
		ParamPageSize:   c.Query(ParamPageSize),
		ParamFilterType: c.Query(ParamFilterType),
		ParamIsClosed:   c.Query(ParamIsClosed),
	}

	// Send Request to caravan service to fetch reports
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchReportsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return
	}

	//Validate response and return if not successfull
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse := apiCR.Response

	// Iterate over reports and check if any of the reports are of type post or comment
	if reports, ok := dataResponse["reports"]; ok {

		// Get Posts & comments data for the reports
		posts, comments, users := fetchReportsEntityData(c, userId, reports.([]interface{}))

		// Add data to response if not empty
		if posts != nil {
			dataResponse["posts"] = posts
		}
		if comments != nil {
			dataResponse["comments"] = comments
		}
		if users != nil {
			dataResponse["users"] = users
		}

	}

	//Generate response
	utils.GenerateResponse(c, dataResponse)
}

func parsePushReportRequest(c *gin.Context) (*PushReportRequest, error) {
	//POST body params
	var prr PushReportRequest

	if err := c.ShouldBindJSON(&prr); err != nil {
		return nil, err
	}

	if prr.CollabcardId != nil {
		prr.CollabcardId = utils.ParseInterfaceToString(prr.CollabcardId)
	}

	if prr.ConversationID != nil {
		prr.ConversationID = utils.ParseInterfaceToString(prr.ConversationID)
	}

	// If chatroom_id is present, parse it & set it to collabcard_id
	if prr.ChatroomId != nil {
		prr.CollabcardId = utils.ParseInterfaceToString(prr.ChatroomId)
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

// Internal method to fetch posts and comments data for the reports
func fetchReportsEntityData(c *gin.Context, userId string, reports []interface{}) (map[string]interface{}, map[string]interface{}, map[string]user.MemberMeta) {

	var post_ids []string
	var comment_ids []string
	var user_ids []string

	var posts map[string]interface{}
	var comments map[string]interface{}
	var users map[string]user.MemberMeta

	// Iterate over reports and get post and comment ids
	for _, report := range reports {

		typeValue, ok := report.(map[string]interface{})["type"]

		if ok {
			if int(typeValue.(float64)) == feed.POST_REPORT_TYPE {
				post_ids = append(post_ids, report.(map[string]interface{})["entity_id"].(string))
			}

			if int(typeValue.(float64)) == feed.COMMENT_REPORT_TYPE || int(typeValue.(float64)) == feed.REPLY_REPORT_TYPE {
				comment_ids = append(comment_ids, report.(map[string]interface{})["entity_id"].(string))
			}

		}
	}

	// if comment_ids are not empty then fetch comments data
	if len(comment_ids) > 0 {

		// create params for the request
		params := map[string]string{
			feed.ParamCommentIds: utils.ParseStringArrayToString(comment_ids),
			feed.ParamUserIsCm:   "true",
		}

		//Send Request to swarm service
		respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.FetchCommentsEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes != nil {

			//Validate and parse response
			response := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)
			if response != nil {
				comments = response["comments"].(map[string]interface{})

				// Get user_ids and post_ids from comments
				for _, comment := range comments {
					user_ids = append(user_ids, comment.(map[string]interface{})["user_id"].(string))

					// If comment is reply then get parent comment's user id
					if parentComment, ok := comment.(map[string]interface{})["parent_comment"]; ok {
						if parentComment != nil {
							user_ids = append(user_ids, parentComment.(map[string]interface{})["user_id"].(string))
						}
					}

					// Get post_id from comments
					post_ids = append(post_ids, comment.(map[string]interface{})["post_id"].(string))
				}

			}
		}
	}

	// if post_ids are not empty then fetch posts data
	if len(post_ids) > 0 {

		// create params for the request
		params := map[string]string{
			feed.ParamPostIds:  utils.ParseStringArrayToString(post_ids),
			feed.ParamUserIsCm: "true",
		}

		//Send request to swarm service
		respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.FetchPostsEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

		//Validate and parse response
		response := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)
		if response != nil {
			posts = response["posts"].(map[string]interface{})

			// Iterate over posts and get user ids
			for _, post := range posts {
				user_ids = append(user_ids, post.(map[string]interface{})["user_id"].(string))
			}
		}

	}

	// If user_ids are not empty, get users data
	if len(user_ids) > 0 {
		var err error

		// Call Internal method to fetch users data
		users, err = user.FetchMemberMeta(utils.CreateHeaders(c, userId), user_ids)
		if err != nil {
			log.Println("Error while fetching users data for reports:", err)
		}
	}

	return posts, comments, users

}
