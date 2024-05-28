package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/handlers/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UpdateUserTopicsRequest struct {
	TopicsIds map[string]bool `json:"topic_ids"`
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

	// Check if user topics are enabled
	if !utils.UserTopicsConnectionEnabled(utils.GetRedisClientFromContext(c), utils.CreateHeaders(c, userId)) {
		utils.GeneralBadRequestError(c, utils.ErrorUserTopicsSettingsNotEnabled)
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
		fetchUsersTopicsInternal(c, userId, response.IsCm)
	case utils.PatchMethod:
		updateUserTopicsInternal(c, userId, response.IsCm)
	}
}

// Internal method to fetch users topics
func fetchUsersTopicsInternal(c *gin.Context, userId string, isCm bool) {

	headers := utils.CreateHeaders(c, userId)

	uuids := c.Query(ParamUUIDs)
	if uuids == "" {
		utils.GeneralBadRequestError(c, "Please send UUIDs in param")
		return
	}

	// Fetch user_unique_ids from uuids
	userIds, err := utility.FetchUserUniqueIdsFromAnyUserIds(headers, uuids)
	if err != nil || len(userIds) == 0 {
		utils.GeneralBadRequestError(c, "Invalid UUIDs sent")
		return
	}

	params := map[string]string{
		ParamUUIDs: utils.ParseStringArrayToString(userIds),
	}

	//add CM role in headers if user is cm
	if isCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	// send request to fetch user topics
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchUserTopicsEndPoint, utils.GETRequest, headers, params, nil)

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse := apiCR.Response

	// fetch user meta for userIds
	userMetaMap, err := utils.FetchMemberMetaMapForUserUniqueIds(utils.GetRedisClientFromContext(c), headers, userIds)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	dataResponse["users"] = userMetaMap

	utils.GenerateResponse(c, dataResponse, true)
}

// Internal method to update user topics
func updateUserTopicsInternal(c *gin.Context, userId string, isCm bool) {

	paramUUID := c.Param("uuid")
	paramUserUniqueId, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), paramUUID)
	if err != nil || paramUserUniqueId == "" {
		utils.GeneralBadRequestError(c, "Please send a valid UUID")
		return
	}

	if !isCm && (paramUserUniqueId != userId) {
		utils.MemberAccessFailError(c)
		return
	}

	// Parse request body
	updateUserTopicsRequest, err := parseUpdateUserTopicsRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//add CM role in headers if user is cm
	if isCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	endpoint := fmt.Sprintf(UpdateUserTopicsEndPoint, paramUserUniqueId)

	// Send request to update user topics
	utils.SendRequest(c, utils.SwarmService, endpoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, updateUserTopicsRequest)
}
