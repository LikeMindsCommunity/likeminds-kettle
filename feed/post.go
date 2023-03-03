package feed

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type OGTags struct {
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type AttachmentMeta struct {
	Name      string `json:"name"`
	Url       string `json:"url"`
	Format    string `json:"format"`
	Size      int    `json:"size"`
	Duration  int    `json:"duration"`
	PageCount int    `json:"page_count"`
	OgTags    OGTags `json:"og_tags"`
}

type AttachmentRequest struct {
	AttachmentType int            `json:"attachment_type" binding:"required"`
	AttachmentMeta AttachmentMeta `json:"attachment_meta"`
}

type CreatePostRequest struct {
	Text        string              `json:"text"`
	Heading     string              `json:"heading"`
	Attachments []AttachmentRequest `json:"attachments"`
	FeedroomID  int                 `json:"feedroom_id"`
}

type DeletePostRequest struct {
	DeleteReason string `json:"delete_reason"`
	UserIsCm     bool   `json:"user_is_cm"`
}

func parseCreatePostRequest(c *gin.Context) (*CreatePostRequest, error) {
	//POST body params
	var cpr CreatePostRequest

	if err := c.ShouldBindJSON(&cpr); err != nil {
		return nil, err
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
		utils.GenerateResponse(c, post_data)

	case utils.POSTMethod:
		createPostInternal(c, userId)

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
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION)
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
	if value, ok := dataResponse["post"]; ok {
		post_data := value.(map[string]interface{})
		user_ids := []string{}

		//Fetch post user id
		if post_user_unique_id, ok := post_data["user_id"]; ok {
			user_ids = append(user_ids, post_user_unique_id.(string))
		}

		//Fetch replies user id
		if replies, ok := post_data["replies"]; ok {
			for _, reply_data := range replies.([]interface{}) {
				if user_unique_id, ok := reply_data.(map[string]interface{})["user_id"]; ok {
					user_ids = append(user_ids, user_unique_id.(string))
				}
			}
		}

		//Fetch user data for given user_unique_ids
		success, user_data := user.FetchMemberMeta(c, user_ids)
		if !success {
			return nil
		}

		//Validation of post based on community member
		if post_user_unique_id, ok := post_data["user_id"]; ok {
			if post_user, ok := user_data[post_user_unique_id.(string)]; ok {
				if post_user.IsDeleted {
					utils.GeneralBadRequestError(c, "Invalid post_id sent!")
					return nil
				}
			}
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data
	}

	return dataResponse
}

func createPostInternal(c *gin.Context, userId string) {
	//Body to be sent in the /post POST request
	createPostRequest, err := parseCreatePostRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, CREATE_POST_ACTION)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, CreatePostEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createPostRequest)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	if createPostRequest.FeedroomID != 0 {
		//Params to be sent in the api/collabcard_follow request
		params := map[string]string{
			chatroom.ParamCollabcardId: strconv.Itoa(createPostRequest.FeedroomID),
			chatroom.ParamMemberId:     userId,
			chatroom.ParamValue:        "true",
		}

		//Send Request to follow the chatroom for the post creator
		utils.SendRequest(c, utils.CoreService, chatroom.CollabcardFollowEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	} else {
		//Generate Response
		utils.GenerateResponse(c, apiCR.Response)
	}
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
	post_user_unique_id := post_data["user_id"]

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
		success, response := user.FetchMemberAccess(c, DELETE_POST_ACTION)
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
