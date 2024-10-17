package feed

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type CreatePostCommentRequest struct {
	TempID      *string                   `json:"temp_id"`
	Text        string                    `json:"text"`
	Attachments []utils.AttachmentRequest `json:"attachments"`
	UUIDs       []string                  `json:"uuids"`
	CreatedAt   int                       `json:"created_at"`
}
type EditCommentRequest struct {
	Text        string                    `json:"text"`
	Attachments []utils.AttachmentRequest `json:"attachments"`
	UserIsCm    bool                      `json:"user_is_cm"`
}

func parseCreateCommentRequest(c *gin.Context) (*CreatePostCommentRequest, error) {
	//POST body params
	var cpcr CreatePostCommentRequest
	raw_data, _ := c.GetRawData()

	if err := json.Unmarshal(raw_data, &cpcr); err != nil {
		return nil, err
	}

	// Iterate over attachments and add widgetsData to widget_meta
	if cpcr.Attachments != nil {
		cpcr.Attachments = utils.ConvertAttachmentMetaForCustomWidgetAttachments(cpcr.Attachments, raw_data)
	}

	return &cpcr, nil
}
func parseEditCommentRequest(c *gin.Context) (*EditCommentRequest, error) {

	//POST body params
	var cpcr EditCommentRequest
	raw_data, _ := c.GetRawData()

	if err := json.Unmarshal(raw_data, &cpcr); err != nil {
		return nil, err
	}

	// Iterate over attachments and add widgetsData to widget_meta
	if cpcr.Attachments != nil {
		cpcr.Attachments = utils.ConvertAttachmentMetaForCustomWidgetAttachments(cpcr.Attachments, raw_data)
	}

	return &cpcr, nil
}

// CommentPost is used to comment on a post
func CommentPost(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	post_id := c.Param("post_id")
	CommentPostEndPoint := fmt.Sprintf(SinglePostCommentEndPoint, post_id)

	//Body to be sent in the /post/<post_id>/comment POST request
	createPostCommentRequest, err := parseCreateCommentRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
	}

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, CREATE_COMMENT_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	// Add member role to headers
	utils.AddMemberRoleToHeaders(c, response.IsCm)

	//Get tagged users from text
	taggedUsers := getTaggedUsersFromText(utils.CreateHeaders(c, userId), createPostCommentRequest.Text)
	createPostCommentRequest.UUIDs = taggedUsers

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, CommentPostEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createPostCommentRequest)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	dataResponse = populateCommentDataResponse(c, dataResponse)

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}

// EditPostComment is used to edit a comment on a post
func EditCommentPost(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	post_id := c.Param("post_id")
	comment_id := c.Param("comment_id")
	editPostCommentEndPoint := fmt.Sprintf(SingleCommentEndPoint, post_id, comment_id)
	GetCommentEndPoint := fmt.Sprintf(SingleCommentEndPoint, post_id, comment_id)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, GetCommentEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if _, ok := dataResponse["comment"]; !ok {
		utils.GeneralBadRequestError(c, "Invalid comment_id sent!")
		return
	}

	//Fetch comment user id
	comment_data := dataResponse["comment"].(map[string]interface{})
	comment_user_unique_id := comment_data["user_id"]

	//Body to be sent in the /post/<post_id>/comment/<comment_id> PUT request
	editCommentRequest, err := parseEditCommentRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
	}

	//If the user is not the comment creator
	if comment_user_unique_id != userId {
		//Fetch member access to view post
		success, response := user.FetchMemberAccess(c, EDIT_COMMENT_ACTION, userId)
		if !success {
			return
		}

		//If not access
		if !response.Access {
			utils.MemberAccessFailError(c)
			return
		}

		// If user is CM, set user_is_cm to true
		editCommentRequest.UserIsCm = response.IsCm
	}

	//Send Request
	respBytes, statusCode = utils.GetRequestResponse(c, utils.SwarmService, editPostCommentEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editCommentRequest)

	//Validate response
	apiCR = utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse = apiCR.Response
	dataResponse = populateCommentDataResponse(c, dataResponse)

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}
