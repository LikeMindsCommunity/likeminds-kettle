package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type DeleteCommentRequest struct {
	DeleteReason string `json:"delete_reason"`
	UserIsCm     bool   `json:"user_is_cm"`
}

func parseDeleteCommentRequest(c *gin.Context) (*DeleteCommentRequest, error) {
	//POST body params
	var dcr DeleteCommentRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}

// CreateCommentReply is used to create a new reply on a comment
func CreateCommentReply(c *gin.Context) {
	Comment(c, utils.POSTMethod)
}

// GetComment is used to get a specific comment
func GetComment(c *gin.Context) {
	Comment(c, utils.GETMethod)
}

// DeletePost is used to delete an existing comment
func DeleteComment(c *gin.Context) {
	Comment(c, utils.DELETEMethod)
}

// Comment method handles post objects
func Comment(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:
		getCommentInternal(c, userId)

	case utils.POSTMethod:
		createCommentInternal(c, userId)

	case utils.DELETEMethod:
		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		deleteCommentInternal(c, userId)
	}
}

func getCommentInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	post_id := c.Param("post_id")
	comment_id := c.Param("comment_id")
	GetCommentEndPoint := fmt.Sprintf(SingleCommentEndPoint, post_id, comment_id)

	//Params to be sent in the /post/<post_id>/comment/<comment_id> request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Fetch member access to view comment
	success, response := user.FetchMemberAccess(c, VIEW_COMMENT_ACTION)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Param updatiion
	params[ParamUserIsCm] = fmt.Sprint(response.IsCm)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, GetCommentEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["comment"]; ok {
		comment_data := value.(map[string]interface{})
		user_ids := []string{}

		//Fetch comment user id
		if comment_user_unique_id, ok := comment_data["user_id"]; ok {
			user_ids = append(user_ids, comment_user_unique_id.(string))
		}

		//Fetch replies user id
		if replies, ok := comment_data["replies"]; ok {
			for _, reply_data := range replies.([]interface{}) {
				if user_unique_id, ok := reply_data.(map[string]interface{})["user_id"]; ok {
					user_ids = append(user_ids, user_unique_id.(string))
				}
			}
		}

		//Fetch user data for given user_unique_ids
		success, user_data := user.FetchMemberMeta(c, user_ids)
		if !success {
			return
		}

		//Validation of comment based on community member
		if comment_user_unique_id, ok := comment_data["user_id"]; ok {
			if comment_user, ok := user_data[comment_user_unique_id.(string)]; ok {
				if comment_user.IsDeleted {
					utils.GeneralBadRequestError(c, "Invalid comment_id sent!")
					return
				}
			}
		}

		//Update users data in dataResponse
		dataResponse["users"] = user_data
	}

	//Send response
	utils.GenerateResponse(c, dataResponse)
}

func createCommentInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	post_id := c.Param("post_id")
	comment_id := c.Param("comment_id")
	CreateCommentEndPoint := fmt.Sprintf(SingleCommentReplyEndPoint, post_id, comment_id)

	//Body to be sent in the /post/<post_id>/comment/<comment_id> POST request
	createCommentReplyRequest, err := parseCreateCommentRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, CREATE_COMMENT_ACTION)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, CreateCommentEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createCommentReplyRequest)
}

func deleteCommentInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	post_id := c.Param("post_id")
	comment_id := c.Param("comment_id")
	DeleteCommentEndPoint := fmt.Sprintf(SingleCommentEndPoint, post_id, comment_id)

	//Body to be sent in the /post/<post_id> DELETE request
	deleteCommentRequest, err := parseDeleteCommentRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to delete post
	success, response := user.FetchMemberAccess(c, DELETE_COMMENT_ACTION)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Update requests body
	deleteCommentRequest.UserIsCm = response.IsCm

	//Send Request
	utils.SendRequest(c, utils.SwarmService, DeleteCommentEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, deleteCommentRequest)

}
