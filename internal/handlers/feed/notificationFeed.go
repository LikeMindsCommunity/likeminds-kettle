package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type CreateUserActivityRequest struct {
	Action string `json:"action" binding:"required"`
}

func parseCreateUserActivityRequest(c *gin.Context) (*CreateUserActivityRequest, error) {
	//POST body params
	var cuar CreateUserActivityRequest

	if err := c.ShouldBindJSON(&cuar); err != nil {
		return nil, err
	}

	return &cuar, nil
}

// FetchNotificationFeed | get user notification feed from swarm service
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
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, GetUserActivityEndPoint, utils.GETRequest, utils.CreateHeaders(c, userID), params, gin.H{})
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

	utils.SendRequest(c, utils.SwarmService, GetUserActivityUnreadCountEndPoint, utils.GETRequest, utils.CreateHeaders(c, userID), params, gin.H{})
}

// NotificationFeedActivityMarkRead | mark user activity as read to swarm service
func NotificationFeedActivityMarkRead(c *gin.Context) {
	userID := user.GetRequestingUserId(c)
	if userID == "" {
		return
	}

	activityID := c.Param("activity_id")

	endpoint := fmt.Sprintf(UserActivityMarkReadEndPoint, activityID)

	utils.SendRequest(c, utils.SwarmService, endpoint, utils.POSTMethod, utils.CreateHeaders(c, userID), nil, gin.H{})
}

// CreateUserActivity is used to create activity for a user
func CreateUserActivity(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	user_id := c.Param("user_id")

	//Get user_unique_id from user_id internally
	user_id, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), user_id)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	UserActivityEndPoint := fmt.Sprintf(SingleUserActivityEndPoint, user_id)

	//Body to be sent in the /user/<user_id>/activity POST request
	createUserActivityRequest, err := parseCreateUserActivityRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
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
	utils.SendRequest(c, utils.SwarmService, UserActivityEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createUserActivityRequest)
}
