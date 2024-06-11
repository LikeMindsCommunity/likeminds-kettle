package feed

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/requests"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type EditPendingPostRequest struct {
	EditPostRequest
	UUIDs []string `json:"uuids"`
}

// CreatePendingPost is used to create a new pending post
func CreatePendingPost(c *gin.Context) {
	PendingPost(c, utils.POSTMethod)
}

// EditPendingPost is used to edit a pending post
func EditPendingPost(c *gin.Context) {
	PendingPost(c, utils.PUTMethod)
}

// FetchPendingPost is used to fetch a pending post
func FetchPendingPost(c *gin.Context) {
	PendingPost(c, utils.GETMethod)
}

// DeletePendingPost is used to delete a pending post
func DeletePendingPost(c *gin.Context) {
	PendingPost(c, utils.DELETEMethod)
}

// Pending post method handles pending post objects
func PendingPost(c *gin.Context, method int) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.POSTMethod:
		createPendingPostInternal(c, userId)

	case utils.PUTMethod:
		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}
		editPendingPostInternal(c, userId)

	case utils.GETMethod:
		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		fetchPendingPostInternal(c, userId)

	case utils.DELETEMethod:
		additionalHeaders := map[string]string{
			utils.HeadersPlatformType: string(utils.PlatformDashboard),
		}
		utils.AddHeaders(c, additionalHeaders)
		botId := user.GetBotId(c)
		deletePendingPostInternal(c, userId, botId)
	}
}

// Internal method to create a pending post for review
func createPendingPostInternal(c *gin.Context, userId string) {
	// Use Create post body params to create Pending post
	cppr, err := parseCreatePostRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	access, _ := checkAccessForCreatPost(c, userId, isPollInPostAttachments(cppr.Attachments))
	if !access {
		return
	}

	// Send request to "/post/pending" of swarm service
	utils.SendRequest(c, utils.SwarmService, CreatePendingPostEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, cppr)
}

// Internal method to edit a pending post
func editPendingPostInternal(c *gin.Context, userId string) {
	pendingPostId := c.Param(ParamPendingPostId)

	// Fetch member access to check whether member is CM or not
	success, response := user.FetchMemberAccess(c, IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//add Admin role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	editPostEndPoint := fmt.Sprintf(PendingPostEndPoint, pendingPostId)

	// Body to be sent in the /post/pending/<pending_post_id> PUT request
	eppr, err := parseEditPostRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get tagged users from text
	taggedUsers := getTaggedUsersFromText(utils.CreateHeaders(c, userId), eppr.Text)

	editPendingPostRequest := EditPendingPostRequest{
		EditPostRequest: *eppr,
		UUIDs:           taggedUsers,
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, editPostEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editPendingPostRequest)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds populate post data
	dataResponse := apiCR.Response
	dataResponse = populatePendingPostDataResponse(c, dataResponse)

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}

// Internal method to fetch a pending post
func fetchPendingPostInternal(c *gin.Context, userId string) {
	pendingPostId := c.Param("pending_post_id")

	fetchPostEndPoint := fmt.Sprintf(PendingPostEndPoint, pendingPostId)

	// Fetch member access to view topics
	success, response := user.FetchMemberAccess(c, IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//add Admin role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, fetchPostEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds populate post data
	dataResponse := apiCR.Response
	dataResponse = populatePendingPostDataResponse(c, dataResponse)

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}

// Internal method to delete pending post
func deletePendingPostInternal(c *gin.Context, userId string, botId string) {
	pendingPostId := c.Param("pending_post_id")
	deletePostEndPoint := fmt.Sprintf(PendingPostEndPoint, pendingPostId)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, deletePostEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds populate post data
	dataResponse := apiCR.Response

	reportId := int(dataResponse["report_id"].(float64))
	if reportId != 0 {
		crnr := requests.CloseReportsNewRequest{ReportIds: []int{reportId}}

		// Send Request to api/community/report
		respBytes, statusCode = utils.GetRequestResponse(c, utils.CoreService, constants.CommunityReportV1EndPoint, utils.PATCHRequest, utils.CreateHeaders(c, botId), nil, crnr)
		if statusCode != http.StatusOK {
			logging.Error(fmt.Sprintf("Error closing report: %v", string(respBytes)))
		}
	}

	dataResponse = map[string]interface{}{
		"success": true,
	}

	utils.GenerateResponse(c, dataResponse, false)
}

// Internal method to populate users data
func populatePendingPostDataResponse(c *gin.Context, dataResponse map[string]interface{}) map[string]interface{} {
	if value, ok := dataResponse["post"]; ok {
		postData := value.(map[string]interface{})
		userIds := []string{}

		// Fetch post user id
		if post_user_unique_id, ok := postData["uuid"]; ok {
			userIds = append(userIds, post_user_unique_id.(string))
		}

		// Get userId
		userId := user.GetRequestingUserId(c)
		redisClient := utils.GetRedisClientFromContext(c)
		headers := utils.CreateHeaders(c, userId)

		//Fetch user data for given user_unique_ids
		user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, userIds)
		if err != nil {
			utils.GeneralAPIError(c, utils.ErrorFetchingUserData)
			return nil
		}

		//Validation of post based on community member
		if post_user_unique_id, ok := postData["uuid"]; ok {
			if post_user, ok := user_data[post_user_unique_id.(string)]; ok && post_user.IsDeleted {
				utils.GeneralBadRequestError(c, utils.ErrorInvalidPostIdSent)
				return nil
			}
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userIds)
	}

	return dataResponse
}
