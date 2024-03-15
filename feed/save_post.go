package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

// SavePost is used to save a post
func CreateSavePost(c *gin.Context) {
	SavePost(c, utils.PUTMethod)
}

// GetSavedPosts is used to get saved posts of a specific user
func GetSavedPosts(c *gin.Context) {
	SavePost(c, utils.GETMethod)
}

// SavePost method handles post save objects
func SavePost(c *gin.Context, method int) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:
		getSavePostsInternal(c, userId)

	case utils.PUTMethod:
		createSavePostInternal(c, userId)

	}
}

func getSavePostsInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	user_id := c.Param("user_id")

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent in the /user/<user_id>/save request
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

	//Param updatiion
	params[ParamUserIsCm] = fmt.Sprint(response.IsCm)

	//Get user_unique_id from user_id internally
	user_id, err := utility.GetUUIDInternally(headers, user_id)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	FetchSavePostEndPoint := fmt.Sprintf(FetchUserSavedPostsEndPoint, user_id)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchSavePostEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["posts"]; ok {

		posts := value.([]interface{})

		redisClient := utils.GetRedisClientFromContext(c)

		user_data, userUniqueIds, err := utils.GetUsersMetaFromFeedData(redisClient, headers, posts, dataResponse)
		if err != nil {
			utils.GenerateResponse(c, nil, false)
			return
		}

		//Update user data in dataResponse
		dataResponse["users"] = user_data

		// Update user topics data in dataResponse
		dataResponse = utils.FetchAndUpdateUserTopicsDataForResponse(redisClient, headers, dataResponse, userUniqueIds)
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

func createSavePostInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	post_id := c.Param("post_id")
	SavePostEndPoint := fmt.Sprintf(SinglePostSaveEndPoint, post_id)

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, SAVE_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, SavePostEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, nil)
}
