package feed

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type UserCommentsAPIResponse struct {
	Success  bool                              `json:"success"`
	Comments []map[string]interface{}          `json:"comments"`
	Posts    map[string]map[string]interface{} `json:"posts"`
	Response map[string]interface{}            `json:"-"`
}

func GetUserComments(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

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
	paramUUID := c.Param("uuid")
	if paramUUID == "" {
		utils.GeneralBadRequestError(c, utils.ErrorInvalidUserId)
	}

	//Get user_unique_id from user_id internally
	userUUID, err := utility.GetUUIDInternally(headers, paramUUID)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	endpoint := fmt.Sprintf(UserCommentsEndPoint, userUUID)

	// Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endpoint, utils.GETRequest, headers, params, nil)

	var userCommentsAPIResponse UserCommentsAPIResponse
	var dataResponse map[string]interface{}

	if statusCode != http.StatusOK {
		// Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return
		}

	} else {
		if err := parseUserCommentsResponse(respBytes, &userCommentsAPIResponse); err != nil {
			utils.GeneralBadRequestError(c, err.Error())
		}

		dataResponse = userCommentsAPIResponse.Response
	}

	// Fetch user data for given user_unique_ids
	userIds := []string{userUUID}

	user_data, userIds, err := populateUserData(c, userId, userCommentsAPIResponse.Comments, userCommentsAPIResponse.Posts, userIds)
	if err != nil {
		utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
		return
	}

	dataResponse["users"] = user_data

	redisClient := utils.GetRedisClientFromContext(c)

	// Update user topics data in dataResponse
	dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userIds)

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

func parseUserCommentsResponse(respBytes []byte, ucar *UserCommentsAPIResponse) error {
	if err := json.Unmarshal(respBytes, &ucar); err != nil {
		return err
	}

	if err := json.Unmarshal(respBytes, &ucar.Response); err != nil {
		return err
	}

	delete(ucar.Response, "success")
	delete(ucar.Response, "error_message")

	return nil
}

func populateUserData(c *gin.Context, userId string, listData []map[string]interface{}, mapData map[string]map[string]interface{}, userIds []string,
) (map[string]utils.MemberMeta, []string, error) {
	userIdsMap := map[string]interface{}{}

	if userIds == nil {
		userIds = []string{}
	}

	for _, userId := range userIds {
		userIdsMap[userId] = nil
	}

	// Fetch user ids from []map[string]interface{} type data
	for _, data := range listData {
		if user_unique_id, ok := data["uuid"]; ok {
			if _, ok := userIdsMap[user_unique_id.(string)]; !ok {
				userIds = append(userIds, user_unique_id.(string))
				userIdsMap[user_unique_id.(string)] = nil
			}
		}
	}

	// Fetch user ids from map[string]map[string]interface{} type data
	for _, data := range mapData {
		if user_unique_id, ok := data["uuid"]; ok {
			if _, ok := userIdsMap[user_unique_id.(string)]; !ok {
				userIds = append(userIds, user_unique_id.(string))
				userIdsMap[user_unique_id.(string)] = nil
			}
		}
	}

	user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(utils.GetRedisClientFromContext(c), utils.CreateHeaders(c, userId), userIds)
	if err != nil {
		return nil, nil, err
	}

	return user_data, userIds, nil
}
