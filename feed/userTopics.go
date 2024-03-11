package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UpdateUserTopicsRequest struct {
	TopicsIds map[string]bool `json:"topics_ids"`
	UserIsCm  bool            `json:"user_is_cm"`
}

func parseUpdateUserTopicsRequest(c *gin.Context) (*UpdateUserTopicsRequest, error) {
	// PATCH body params
	var uutr UpdateUserTopicsRequest

	if err := c.ShouldBindJSON(&uutr); err != nil {
		return nil, err
	}

	return &uutr, nil
}

// Exposed method to fetch users topics
func FetchUsersTopics(c *gin.Context) {
	userTopics(c, utils.GETMethod)
}

// Exposed method to update user topics
func UpdateUserTopics(c *gin.Context) {
	userTopics(c, utils.PatchMethod)
}

func userTopics(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Fetch member access if user is a member
	success, response := user.FetchMemberAccess(c, IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	switch method {
	case utils.GETMethod:
		fetchUsersTopicsInternal(c, userId)
	case utils.PatchMethod:
		updateUserTopicsInternal(c, userId, response.IsCm)
	}
}

// Internal method to fetch users topics
func fetchUsersTopicsInternal(c *gin.Context, userId string) {

	uuids := utils.ParseStringArrayFromParam(c.Query(ParamUUIDs))
	if len(uuids) <= 0 {
		utils.GeneralBadRequestError(c, "Please send a valid UUID in params")
	}

	interfaceUUIDs := make([]interface{}, len(uuids))
	for i, v := range uuids {
		interfaceUUIDs[i] = v
	}

	// Fetch user_ids from uuids
	userIds, err := utility.GetUsersInfoInternally(utils.CreateHeaders(c, userId), interfaceUUIDs, true)
	if err != nil {
		return
	}

	params := map[string]string{
		ParamUUIDs: utils.ParseArrayToString(userIds.([]interface{})),
	}

	// send request to fetch user topics
	utils.SendRequest(c, utils.SwarmService, FetchUserTopicsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}

// Internal method to update user topics
func updateUserTopicsInternal(c *gin.Context, userId string, isCm bool) {

	paramUUID := c.Param("user_id")
	paramUserId, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), paramUUID)
	if err != nil || paramUserId == "" {
		utils.GeneralBadRequestError(c, "Please send a valid UUID")
		return
	}

	if !isCm && (paramUserId != userId) {
		utils.MemberAccessFailError(c)
		return
	}

	// Parse request body
	updateUserTopicsRequest, err := parseUpdateUserTopicsRequest(c)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	updateUserTopicsRequest.UserIsCm = isCm

	endpoint := fmt.Sprintf(UpdateUserTopicsEndPoint, paramUserId)

	// Send request to update user topics
	utils.SendRequest(c, utils.SwarmService, endpoint, utils.PatchMethod, utils.CreateHeaders(c, userId), nil, updateUserTopicsRequest)
}
