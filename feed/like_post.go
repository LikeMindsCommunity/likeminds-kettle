package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateLikeRequest struct {
	CreatedAt int `json:"created_at"`
}

func parseCreateLikeRequest(c *gin.Context) (*CreateLikeRequest, error) {
	//POST body params
	var clr CreateLikeRequest

	if err := c.ShouldBindJSON(&clr); err != nil {
		return nil, err
	}

	return &clr, nil
}

// CreatePostLike is used to like on a post
func CreatePostLike(c *gin.Context) {
	PostLike(c, utils.PUTMethod)
}

// GetPostLikes is used to get likes of a specific post
func GetPostLikes(c *gin.Context) {
	PostLike(c, utils.GETMethod)
}

// PostLike method handles post like objects
func PostLike(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	post_id := c.Param("post_id")
	LikePostEndPoint := fmt.Sprintf(SinglePostLikeEndPoint, post_id)

	//Send request
	switch method {
	case utils.GETMethod:
		getPostLikesInternal(c, userId, LikePostEndPoint)

	case utils.PUTMethod:
		createPostLikeInternal(c, userId, LikePostEndPoint)

	}
}

func getPostLikesInternal(c *gin.Context, userId string, endPoint string) {
	//Params to be sent in the /post/<post_id>/like request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Fetch member access to view post likes
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION, userId)
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
		user_data, err := utils.FetchMemberMetaMapFromCache(utils.GetRedisClientFromContext(c), utils.CreateHeaders(c, userId), user_ids)
		if err != nil {
			utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
			return
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

func createPostLikeInternal(c *gin.Context, userId string, endPoint string) {
	//Body to be sent in the /post/<post_id>/like PUT request
	createPostLikeRequest, _ := parseCreateLikeRequest(c)

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, LIKE_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, endPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, createPostLikeRequest)
}
