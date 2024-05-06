package feed

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreatePostRequest struct {
	TempID         *string                   `json:"temp_id"`
	TopicIDs       []string                  `json:"topic_ids"`
	Text           string                    `json:"text"`
	Heading        string                    `json:"heading"`
	Attachments    []utils.AttachmentRequest `json:"attachments"`
	FeedroomID     int                       `json:"feedroom_id"`
	UUIDs          []string                  `json:"uuids"`
	OnBehalfOfUUID string                    `json:"on_behalf_of_uuid,omitempty"`
	IsRepost       bool                      `json:"is_repost"`
	Visibility     string                    `json:"visibility,omitempty"`
	UserIsCm       bool                      `json:"user_is_cm,omitempty"`
	CreatedAt      int                       `json:"created_at"`
}

type EditPostRequest struct {
	Text        string                    `json:"text"`
	TopicIDs    []string                  `json:"topic_ids"`
	Heading     string                    `json:"heading,omitempty"`
	Attachments []utils.AttachmentRequest `json:"attachments"`
	Visibility  string                    `json:"visibility,omitempty"`
	UserIsCm    bool                      `json:"user_is_cm"`
}

type DeletePostRequest struct {
	DeleteReason string `json:"delete_reason"`
	UserIsCm     bool   `json:"user_is_cm"`
}

func parseCreatePostRequest(c *gin.Context) (*CreatePostRequest, error) {
	//POST body params
	var cpr CreatePostRequest
	raw_data, _ := c.GetRawData()

	if err := json.Unmarshal(raw_data, &cpr); err != nil {
		return nil, err
	}

	// Unmarshal widgets data for attachment type custom widget
	widgets_data := make(map[string]interface{})
	err := json.Unmarshal(raw_data, &widgets_data)
	if err != nil {
		return nil, err
	}

	// Iterate over attachments and add widgets_data to widget_meta
	if cpr.Attachments != nil {
		cpr.Attachments = utils.ConvertAttachmentMetaForCustomWidgetAttachments(cpr.Attachments, raw_data)
	}

	return &cpr, nil
}

func parseEditPostRequest(c *gin.Context) (*EditPostRequest, error) {
	//POST body params
	var cpr EditPostRequest

	raw_data, _ := c.GetRawData()

	if err := json.Unmarshal(raw_data, &cpr); err != nil {
		return nil, err
	}

	// Iterate over attachments and add widgets_data to widget_meta
	if cpr.Attachments != nil {
		cpr.Attachments = utils.ConvertAttachmentMetaForCustomWidgetAttachments(cpr.Attachments, raw_data)
	}

	return &cpr, nil
}

func parseDeletePostRequest(c *gin.Context) (*DeletePostRequest, error) {
	//POST body params
	var dpr DeletePostRequest

	if err := c.ShouldBindJSON(&dpr); err != nil {
		return nil, err
	}

	return &dpr, nil
}

func populatePostDataResponse(c *gin.Context, dataResponse map[string]interface{}) map[string]interface{} {
	if value, ok := dataResponse["post"]; ok {
		post_data := value.(map[string]interface{})
		user_ids := []string{}

		//Fetch post user id
		if post_user_unique_id, ok := post_data["uuid"]; ok {
			user_ids = append(user_ids, post_user_unique_id.(string))
		}

		//Fetch replies user id
		if replies, ok := post_data["replies"]; ok {
			for _, reply_data := range replies.([]interface{}) {
				if user_unique_id, ok := reply_data.(map[string]interface{})["uuid"]; ok {
					user_ids = append(user_ids, user_unique_id.(string))
				}
			}
		}

		user_ids = utils.AppendRepostPostUsersFromFeedDataResponse(dataResponse, user_ids)
		user_ids = utils.AppendPollOptionAddedByUsersFromFeedDataResponse(dataResponse, user_ids)

		// Get userId
		userId := user.GetRequestingUserId(c)
		redisClient := utils.GetRedisClientFromContext(c)
		headers := utils.CreateHeaders(c, userId)

		//Fetch user data for given user_unique_ids
		user_data, err := utils.FetchMemberMetaMapForUserUniqueIds(redisClient, headers, user_ids)
		if err != nil {
			utils.GeneralAPIError(c, utils.ErrorFetchingUserData)
			return nil
		}

		//Validation of post based on community member
		if post_user_unique_id, ok := post_data["uuid"]; ok {
			if post_user, ok := user_data[post_user_unique_id.(string)]; ok {
				if post_user.IsDeleted {
					utils.GeneralBadRequestError(c, "Invalid post_id sent!")
					return nil
				}
			}
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, user_ids)
	}

	return dataResponse
}

// CreatePost is used to create a new post
func CreatePost(c *gin.Context) {
	Post(c, utils.POSTMethod)
}

// GetPost is used to get a specific post
func GetPost(c *gin.Context) {
	Post(c, utils.GETMethod)
}

// DeletePost is used to delete an existing post
func DeletePost(c *gin.Context) {
	Post(c, utils.DELETEMethod)
}

// EditPost is used to edit an existing post
func EditPost(c *gin.Context) {
	Post(c, utils.PUTMethod)
}

// Post method handles post objects
func Post(c *gin.Context, method int) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:
		postId := c.Param("post_id")
		post_data := GetPostInternal(c, userId, postId)
		if post_data == nil {
			return
		}

		//Send response
		utils.GenerateResponse(c, post_data, true)

	case utils.POSTMethod:
		createPostInternal(c, userId)

	case utils.PUTMethod:
		editPostInternal(c, userId)

	case utils.DELETEMethod:
		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		deletePostInternal(c, userId)
	}
}

func GetPostInternal(c *gin.Context, userId string, postId string) map[string]interface{} {
	//Url generation
	GetPostEndPoint := fmt.Sprintf(SinglePostEndPoint, postId)

	//Params to be sent in the /post/<post_id> GET request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION, userId)
	if !success {
		return nil
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return nil
	}

	//Param updation
	params[ParamUserIsCm] = fmt.Sprint(response.IsCm)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, GetPostEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return nil
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	dataResponse = populatePostDataResponse(c, dataResponse)

	return dataResponse
}

func createPostInternal(c *gin.Context, userId string) {

	headers := utils.CreateHeaders(c, userId)

	//Body to be sent in the /post POST request
	createPostRequest, err := parseCreatePostRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, CREATE_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Update user_is_cm in request
	createPostRequest.UserIsCm = response.IsCm

	if createPostRequest.IsRepost {
		if !utils.FeedRepostSettingsEnabled(utils.GetRedisClientFromContext(c), headers) {
			utils.GeneralBadRequestError(c, utils.ErrorRepostSettingNotEnabled)
			return
		}
	}

	if createPostRequest.OnBehalfOfUUID != "" {
		//Fetch member access to change author
		success, response := user.FetchMemberAccess(c, CHANGE_AUTHOR_ACTION, userId)
		if !success {
			return
		}

		//If not access
		if !response.Access {
			utils.MemberAccessFailError(c)
			return
		}

		//Update user_is_cm in request
		createPostRequest.UserIsCm = response.IsCm

		//Get parsed UUID from core service
		parsedUUID, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), createPostRequest.OnBehalfOfUUID)

		//If error in fetching UUID
		if err != nil || parsedUUID == "" {
			utils.GeneralBadRequestError(c, "Invalid on_behalf_of_uuid sent!")
			return
		}

		//Update UUID in request
		createPostRequest.OnBehalfOfUUID = parsedUUID
	}

	//Get tagged users from text
	taggedUsers := getTaggedUsersFromText(utils.CreateHeaders(c, userId), createPostRequest.Text)
	createPostRequest.UUIDs = taggedUsers
	var createPostEndpoint string

	if !utils.IsPostApprovalNeeded(utils.GetRedisClientFromContext(c), headers) {
		createPostEndpoint = CreatePostEndPoint
	} else {
		createPostEndpoint = CreatePendingPostEndPoint
	}

	// Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, createPostEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createPostRequest)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	dataResponse = populatePostDataResponse(c, dataResponse)

	//If flow succeeds
	if createPostRequest.FeedroomID != 0 {
		//Params to be sent in the api/collabcard_follow request
		params := map[string]string{
			chatroom.ParamCollabcardId: strconv.Itoa(createPostRequest.FeedroomID),
			chatroom.ParamMemberId:     userId,
			chatroom.ParamValue:        "true",
		}

		//Send Request to follow the chatroom for the post creator
		utils.GetRequestResponseWithoutContext(utils.CoreService, chatroom.CollabcardFollowEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}

func editPostInternal(c *gin.Context, userId string) {

	post_id := c.Param("post_id")
	EditPostEndPoint := fmt.Sprintf(SinglePostEndPoint, post_id)
	GetPostEndPoint := fmt.Sprintf(SinglePostEndPoint, post_id)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, GetPostEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if _, ok := dataResponse["post"]; !ok {
		utils.GeneralBadRequestError(c, "Invalid post_id sent!")
		return
	}

	//Fetch post user id
	post_data := dataResponse["post"].(map[string]interface{})
	post_user_unique_id := post_data["uuid"]

	//Body to be sent in the /post/<post_id> PUT request
	editPostRequest, err := parseEditPostRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	if post_user_unique_id != userId {

		// Fetch member access to edit post
		success, response := user.FetchMemberAccess(c, EDIT_POST_ACTION, userId)
		if !success {
			return
		}

		// If not access
		if !response.Access {
			utils.MemberAccessFailError(c)
			return
		}

		// Update user_is_cm in request
		editPostRequest.UserIsCm = response.IsCm
	}

	//Send Request
	respBytes, statusCode = utils.GetRequestResponse(c, utils.SwarmService, EditPostEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editPostRequest)

	//Validate response
	apiCR = utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds populate post data
	dataResponse = apiCR.Response
	dataResponse = populatePostDataResponse(c, dataResponse)

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}

func deletePostInternal(c *gin.Context, userId string) {
	post_id := c.Param("post_id")
	DeletePostEndPoint := fmt.Sprintf(SinglePostEndPoint, post_id)
	GetPostEndPoint := fmt.Sprintf(SinglePostEndPoint, post_id)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, GetPostEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if _, ok := dataResponse["post"]; !ok {
		utils.GeneralBadRequestError(c, "Invalid post_id sent!")
		return
	}

	//Fetch post user id
	post_data := dataResponse["post"].(map[string]interface{})
	post_user_unique_id := post_data["uuid"]

	//Body to be sent in the /post/<post_id> DELETE request
	deletePostRequest, err := parseDeletePostRequest(c)
	if err != nil {
		//If body is not present
		if err.Error() == "EOF" {
			deletePostRequest = &DeletePostRequest{}
		} else {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}
	}

	//If the user is not the post creator
	if post_user_unique_id != userId {
		//Fetch member access to delete post
		success, response := user.FetchMemberAccess(c, DELETE_POST_ACTION, userId)
		if !success {
			return
		}

		//If not access
		if !response.Access {
			utils.MemberAccessFailError(c)
			return
		}

		//Update requests body
		deletePostRequest.UserIsCm = response.IsCm
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, DeletePostEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, deletePostRequest)
}

func getTaggedUsersFromText(headers map[string]interface{}, text string) []string {
	taggedUsers, userUniqueIds := []string{}, []string{}

	// Get user unique id from member route using regex
	pattern, _ := regexp.Compile("route://[member member_profile]+/([a-f0-9]{8}-?[a-f0-9]{4}-?4[a-f0-9]{3}-?[89ab][a-f0-9]{3}-?[a-f0-9]{12})")
	allSubstringMatches := pattern.FindAllStringSubmatch(text, -1)

	for _, occurance := range allSubstringMatches {
		taggedUsers = append(taggedUsers, occurance[1])
	}

	// Get client user unique id from user_profile route using regex
	pattern, _ = regexp.Compile(`route:\/\/user_profile\/([\s\S]*?)>>`)
	allSubstringMatches = pattern.FindAllStringSubmatch(text, -1)

	for _, occurance := range allSubstringMatches {
		taggedUsers = append(taggedUsers, occurance[1])
	}

	// Get valid user unique ids by calling internal users meta api
	if len(taggedUsers) > 0 {
		userUniqueIds, _ = utility.FetchUserUniqueIdsFromAnyUserIds(headers, taggedUsers)
	}

	return userUniqueIds
}
