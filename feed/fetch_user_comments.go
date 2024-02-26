package feed

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UserCommentsResponse struct {
	Comments []map[string]interface{}          `json:"comments"`
	Posts    map[string]map[string]interface{} `json:"posts"`
}

func GetUserComments(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Fetch member access to view post
	success, response := user.FetchMemberAccess(c, VIEW_COMMENT_ACTION, userId)
	if !success {
		return
	}

	// If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	// Params to be sent in the general post search api internally
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
		ParamUserIsCm: fmt.Sprint(response.IsCm),
	}

	// Access query params and url generation
	userID := c.Param("user_id")

	if userID == "" {
		utils.GeneralBadRequestError(c, utils.ErrorInvalidUserId)
	}

	//Get user_unique_id from user_id internally
	userUUID, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), userID)

	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	endpoint := fmt.Sprintf(UserCommentsEndPoint, userUUID)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	// If flow succeeds
	dataResponse := apiCR.Response

	var userCommentsResponse *UserCommentsResponse

	convertedDataResponse, _ := json.Marshal(dataResponse)
	json.Unmarshal(convertedDataResponse, &userCommentsResponse)

	//Fetch user data for given user_unique_ids
	userIds := []string{userUUID}
	userIdsMap := map[string]interface{}{
		userUUID: nil,
	}

	// Fetch user ids from comments data
	for _, commentData := range userCommentsResponse.Comments {
		if user_unique_id, ok := commentData["uuid"]; ok {
			if _, ok := userIdsMap[user_unique_id.(string)]; !ok {
				userIds = append(userIds, user_unique_id.(string))
				userIdsMap[user_unique_id.(string)] = nil
			}
		}
	}

	// Fetch user ids from posts data
	for _, postData := range userCommentsResponse.Posts {
		if user_unique_id, ok := postData["uuid"]; ok {
			if _, ok := userIdsMap[user_unique_id.(string)]; !ok {
				userIds = append(userIds, user_unique_id.(string))
				userIdsMap[user_unique_id.(string)] = nil
			}
		}
	}

	user_data, err := user.FetchMemberMeta(utils.CreateHeaders(c, userId), userIds)
	if err != nil {
		utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
		return
	}

	dataResponse["users"] = user_data

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
