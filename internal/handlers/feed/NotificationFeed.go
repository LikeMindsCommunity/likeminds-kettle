package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/requests"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func parseCreateNotificationActivityRequest(c *gin.Context, userId string,
) (*requests.CreateNotificationActivityRequest, error) {

	//POST body params
	var cuar requests.CreateNotificationActivityRequest
	if err := c.ShouldBindJSON(&cuar); err != nil {
		return nil, err
	}

	headers := utils.CreateHeaders(c, userId)

	// Fetch LM user_unique_id against uuids
	UUID, err := utility.GetUUIDInternally(headers, cuar.ActionBy)
	if err != nil {
		return nil, err
	}
	cuar.ActionBy = UUID

	UUIDs, err := utility.FetchUserUniqueIdsFromAnyUserIds(headers, cuar.ActionOn)
	if err != nil {
		return nil, err
	}
	cuar.ActionOn = UUIDs

	return &cuar, nil
}

// CreateActivityForNotificationFeed is used to create activity for a user
func CreateActivityForNotificationFeed(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the /user/<user_id>/activity POST request
	createUserActivityRequest, err := parseCreateNotificationActivityRequest(c, userId)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, CREATE_ACTIVITY_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, constants.NotificationFeedEndpoint, utils.POSTRequestRawBody,
		utils.CreateHeaders(c, userId), nil, createUserActivityRequest)
}

// FetchNotificationFeed | get user Notification feed from swarm service
func FetchNotificationFeed(c *gin.Context) {
	userID := user.GetRequestingUserId(c)
	if userID == "" {
		return
	}

	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	success, response := user.FetchMemberAccess(c, VIEW_USER_ACTIVITY, userID)
	if !success {
		return
	}

	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	utils.AddMemberRoleToHeaders(c, response.IsCm)

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, constants.NotificationFeedEndpoint, utils.GETRequest, utils.CreateHeaders(c, userID), params, gin.H{})
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)
}

// FetchNotificationFeedUnreadCount | get user activity unread count from swarm service
func FetchNotificationFeedUnreadCount(c *gin.Context) {
	userID := user.GetRequestingUserId(c)
	if userID == "" {
		return
	}

	//Params to be sent in the api/chatroom/fetch_all request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	utils.SendRequest(c, utils.SwarmService, constants.NotificationFeedUnreadCountEndPoint, utils.GETRequest, utils.CreateHeaders(c, userID), params, gin.H{})
}

// NotificationActivityMarkRead | mark user activity as read to swarm service
func NotificationActivityMarkRead(c *gin.Context) {
	userID := user.GetRequestingUserId(c)
	if userID == "" {
		return
	}

	activityID := c.Param("activity_id")

	endpoint := fmt.Sprintf(constants.NotificationActivityMarkReadEndPoint, activityID)

	utils.SendRequest(c, utils.SwarmService, endpoint, utils.POSTMethod, utils.CreateHeaders(c, userID), nil, gin.H{})
}
