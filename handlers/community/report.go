package community

import (
	"fmt"
	"reflect"

	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	log "github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/requests"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/feed"
	"github.com/nateshr/likeminds-authentication/handlers/user"
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

type PushReportV1Request struct {
	EntityId    string `json:"entity_id" binding:"required"`
	EntityType  string `json:"entity_type" binding:"required"`
	AccusedUUID string `json:"accused_uuid"`
	TagId       int    `json:"tag_id,omitempty"`
	Reason      string `json:"reason"`
	Link        string `json:"link"`
}

type CloseReportRequest struct {
	ReportID int `json:"report_id" binding:"required"`
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

func parseCloseReportsNewRequest(c *gin.Context) (*requests.CloseReportsNewRequest, error) {
	//POST body params
	var crr requests.CloseReportsNewRequest

	if err := c.ShouldBindJSON(&crr); err != nil {
		return nil, err
	}

	return &crr, nil
}

func parsePushReportV1Request(c *gin.Context) (*PushReportV1Request, error) {
	//POST body params
	var prr PushReportV1Request

	if err := c.ShouldBindJSON(&prr); err != nil {
		return nil, err
	}

	return &prr, nil
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

// UpdateReports is used to close a report in community
func UpdateReports(c *gin.Context) {
	Report(c, utils.PatchMethod)
}

// Report method handles community reports
func Report(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Get Bot Id if request from dashboard
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:

		getReportsInternal(c, userId)

	case utils.POSTMethod:

		apiRevampCheckV1 := utils.ApiRevampV1Check(c)

		if apiRevampCheckV1 {
			// If ApiRevampV1Check is true, call api/community/report with POST method
			pushReportsInternalV1(c, userId)
		} else {
			// Else call api/push_report with POST method
			pushReportsInternalOld(c, userId)
		}

	case utils.PatchMethod:

		updateReportsInternal(c, userId)

	case utils.DELETEMethod:

		//call api/close_report with POST method
		closeReportsInternalOld(c, userId)

	}
}

func getReportsInternal(c *gin.Context, userId string) {

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent with pagination and filter support in API
	params := map[string]string{
		ParamPage:       c.Query(ParamPage),
		ParamPageSize:   c.Query(ParamPageSize),
		ParamFilterType: c.Query(ParamFilterType),
		ParamIsClosed:   c.Query(ParamIsClosed),
	}

	// Send Request to caravan service to fetch reports
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchReportsEndPoint, utils.GETRequest, headers, params, nil)
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

		redisClient := utils.GetRedisClientFromContext(c)

		// Get Posts & comments data for the reports
		userUniqueIds, posts, comments, topics, widgets, users, repostedPosts := fetchReportsEntityData(headers, redisClient, reports.([]interface{}))

		dataResponse["posts"] = posts
		dataResponse["comments"] = comments
		dataResponse[constants.ResponseKeyTopics] = topics
		dataResponse[constants.ResponseKeyWidgets] = widgets
		dataResponse["users"] = users
		dataResponse["reposted_posts"] = repostedPosts

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	//Generate response
	utils.GenerateResponse(c, dataResponse, true)
}

// Internal method to push reports Old
func pushReportsInternalOld(c *gin.Context, userId string) {

	pushReportRequest, err := parsePushReportRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, PushReportEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, pushReportRequest)
}

// Internal method to push reports ApiRevamp v1
func pushReportsInternalV1(c *gin.Context, userId string) {

	pushReportV1Request, err := parsePushReportV1Request(c)
	if err != nil {
		// Throw error if parsing fails
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, constants.CommunityReportV1EndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, pushReportV1Request)
}

// Internal method to close reports Old
func closeReportsInternalOld(c *gin.Context, userId string) {

	//Body to be sent in the close report api internally
	closeReportRequest, err := parseCloseReportRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request to api/close_report
	utils.SendRequest(c, utils.CoreService, CloseReportsEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, closeReportRequest)
}

// Internal method to close reports ApiRevamp v1
func updateReportsInternal(c *gin.Context, userId string) {

	//Body to be sent in the close report api internally
	crnr, err := parseCloseReportsNewRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request to api/community/report/close
	utils.SendRequest(c, utils.CoreService, constants.CommunityReportV1EndPoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, crnr)
}

func fetchReportsEntityData(headers map[string]interface{}, redisClient *redis.Client, reports []interface{},
) ([]string, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]utils.MemberMeta, map[string]interface{}) {

	posts, comments, topics, widgets, repostedPosts := map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}
	users, userIds := map[string]utils.MemberMeta{}, []string{}

	// Get post_ids, comment_ids and pending_post_ids from reports
	postIds, pendingPostIds, commentIds := extractIdsFromReports(reports)

	// if comment_ids are not empty then fetch comments data
	if len(commentIds) > 0 {
		comments, widgets, userIds, postIds = processComments(headers, commentIds, widgets, userIds, postIds)
	}

	// if post_ids are not empty then fetch posts data
	if len(postIds) > 0 || len(pendingPostIds) > 0 {
		posts, topics, widgets, repostedPosts, userIds = processPosts(headers, postIds, pendingPostIds, topics, widgets, userIds)
	}

	// If user_ids are not empty, get users data
	if len(userIds) > 0 {
		users = processUsers(redisClient, headers, userIds)
	}

	return userIds, posts, comments, topics, widgets, users, repostedPosts
}

func extractIdsFromReports(reports []interface{}) ([]string, []string, []string) {

	postIds, pendingPostIds, commentIds := []string{}, []string{}, []string{}

	// Iterate over reports and get post and comment ids
	for _, report := range reports {

		typeValue, ok := report.(map[string]interface{})["type"]

		// If type is string, convert it to report type int
		if reflect.TypeOf(typeValue).Kind() == reflect.String {
			typeValue = float64(ReportTypeStrintToInt(typeValue.(string)))
		}

		if ok {
			if int(typeValue.(float64)) == feed.POST_REPORT_TYPE {
				postIds = append(postIds, report.(map[string]interface{})[constants.ResponseKeyEntityId].(string))
			}

			if int(typeValue.(float64)) == feed.COMMENT_REPORT_TYPE || int(typeValue.(float64)) == feed.REPLY_REPORT_TYPE {
				commentIds = append(commentIds, report.(map[string]interface{})[constants.ResponseKeyEntityId].(string))
			}

			if int(typeValue.(float64)) == feed.PENDING_POST_REPORT_TYPE {
				pendingPostIds = append(pendingPostIds, report.(map[string]interface{})[constants.ResponseKeyEntityId].(string))
			}

		}
	}

	return postIds, pendingPostIds, commentIds
}

func fetchCommentsData(headers map[string]interface{}, commentIds []string,
) (map[string]interface{}, map[string]interface{}) {

	comments, widgets := map[string]interface{}{}, map[string]interface{}{}

	// create params for the request
	params := map[string]string{
		feed.ParamCommentIds: utils.ParseStringArrayToString(commentIds),
		feed.ParamUserIsCm:   "true",
	}

	//Send Request to swarm service
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.FetchCommentsEndpoint, utils.GETRequest, headers, params, nil)
	if respBytes != nil {

		//Validate and parse response
		response := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)
		if response != nil {
			comments = response["comments"].(map[string]interface{})
			widgets = response[constants.ResponseKeyWidgets].(map[string]interface{})
		}
	}

	return comments, widgets
}

func fetchPostsData(headers map[string]interface{}, postIds []string, pendingPostIds []string) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}) {

	posts, topics, widgets, repostedPosts := map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}

	// create params for the request
	params := map[string]string{
		feed.ParamUserIsCm: "true",
	}

	// If post_ids are not empty, add post_ids to params
	if len(postIds) > 0 {
		params[feed.ParamPostIds] = utils.ParseStringArrayToString(postIds)
	}

	// If pending_post_ids are not empty, add pending_post_ids to params
	if len(pendingPostIds) > 0 {
		params[feed.ParamPendingPostIds] = utils.ParseStringArrayToString(pendingPostIds)
	}

	//Send request to swarm service
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.FetchPostsEndpoint, utils.GETRequest, headers, params, nil)

	//Validate and parse response
	response := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)
	if response != nil {
		posts = response["posts"].(map[string]interface{})
		topics = response["topics"].(map[string]interface{})
		widgets = response[constants.ResponseKeyWidgets].(map[string]interface{})
		repostedPosts = response["reposted_posts"].(map[string]interface{})

	}

	return posts, topics, widgets, repostedPosts
}

func processComments(headers map[string]interface{}, commentIds []string, widgets map[string]interface{}, userIds []string, postIds []string,
) (map[string]interface{}, map[string]interface{}, []string, []string) {

	comments, commentWidgets := fetchCommentsData(headers, commentIds)

	if len(commentWidgets) > 0 {
		for id, widget := range commentWidgets {
			widgets[id] = widget
		}
	}

	// Get user_ids and post_ids from comments
	for _, comment := range comments {
		userIds = append(userIds, comment.(map[string]interface{})["uuid"].(string))

		// If comment is reply then get parent comment's user id
		if parentComment, ok := comment.(map[string]interface{})["parent_comment"]; ok && parentComment != nil {
			userIds = append(userIds, parentComment.(map[string]interface{})["uuid"].(string))
		}

		// Get post_id from comments
		postIds = append(postIds, comment.(map[string]interface{})["post_id"].(string))
	}

	return comments, widgets, userIds, postIds
}

func processPosts(headers map[string]interface{}, postIds []string, pendingPostIds []string, topics map[string]interface{},
	widgets map[string]interface{}, userIds []string,
) (map[string]interface{}, map[string]interface{}, map[string]interface{}, map[string]interface{}, []string) {

	posts, postTopics, postWidgets, repostedPosts := fetchPostsData(headers, postIds, pendingPostIds)

	if len(postWidgets) > 0 {
		for id, widget := range postWidgets {
			widgets[id] = widget
		}
	}

	if len(postTopics) > 0 {
		for id, topic := range postTopics {
			topics[id] = topic
		}
	}

	// Iterate over posts and get user ids
	for _, post := range posts {
		userIds = append(userIds, post.(map[string]interface{})["uuid"].(string))
	}

	for _, repostedPost := range repostedPosts {
		userIds = append(userIds, repostedPost.(map[string]interface{})["uuid"].(string))
	}

	return posts, topics, widgets, repostedPosts, userIds
}

func processUsers(redisClient *redis.Client, headers map[string]interface{}, userIds []string) map[string]utils.MemberMeta {
	users := map[string]utils.MemberMeta{}
	usersMeta, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, userIds)
	if err != nil {
		log.Error(fmt.Sprintf("Error while fetching users data for reports: %s", err))
	}

	for id, userMeta := range usersMeta {
		users[id] = userMeta
	}

	return users
}
