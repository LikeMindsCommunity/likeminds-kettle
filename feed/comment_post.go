package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreatePostCommentRequest struct {
	Text string `json:"text" binding:"required"`
}

func parseCreateCommentRequest(c *gin.Context) (*CreatePostCommentRequest, error) {
	//POST body params
	var cpcr CreatePostCommentRequest

	if err := c.ShouldBindJSON(&cpcr); err != nil {
		return nil, err
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
		utils.GeneralAPIError(c, err.Error())
	}

	//Fetch member access to view post
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
	utils.SendRequest(c, utils.SwarmService, CommentPostEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createPostCommentRequest)
}
