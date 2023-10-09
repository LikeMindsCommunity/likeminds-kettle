package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// CreateCommentLike is used to like on a comment
func CreateCommentLike(c *gin.Context) {
	CommentLike(c, utils.PUTMethod)
}

// GetCommentLikes is used to get likes of a specific comment
func GetCommentLikes(c *gin.Context) {
	CommentLike(c, utils.GETMethod)
}

// CommentLike method handles comment like objects
func CommentLike(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	post_id := c.Param("post_id")
	comment_id := c.Param("comment_id")
	LikeCommentEndPoint := fmt.Sprintf(SingleCommentLikeEndPoint, post_id, comment_id)

	//Send request
	switch method {
	case utils.GETMethod:
		getCommentLikesInternal(c, userId, LikeCommentEndPoint)

	case utils.PUTMethod:
		createCommentLikeInternal(c, userId, LikeCommentEndPoint)

	}
}

func getCommentLikesInternal(c *gin.Context, userId string, endPoint string) {
	//Params to be sent in the /post/<post_id>/commnet/<comment_id>like request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Fetch member access to view post likes
	success, response := user.FetchMemberAccess(c, VIEW_COMMENT_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["likes"]; ok {
		likes_data := value.([]interface{})
		user_ids := []string{}

		//Fetch user ids
		for _, like_data := range likes_data {
			if user_unique_id, ok := like_data.(map[string]interface{})["uuid"]; ok {
				user_ids = append(user_ids, user_unique_id.(string))
			}
		}

		//Fetch user data for given user_unique_ids
		user_data, err := user.FetchMemberMeta(utils.CreateHeaders(c, userId), user_ids)
		if err != nil {
			utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
			return
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data
	}

	//Send response
	utils.GenerateResponse(c, dataResponse)
}

func createCommentLikeInternal(c *gin.Context, userId string, endPoint string) {
	//Body to be sent in the /post/<post_id>/comment/<comment_id>/like PUT request
	createCommentLikeRequest, _ := parseCreateLikeRequest(c)

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, LIKE_COMMENT_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, endPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, createCommentLikeRequest)
}
