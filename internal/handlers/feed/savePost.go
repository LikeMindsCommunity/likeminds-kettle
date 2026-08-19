package feed

import (
	"fmt"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/utility"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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
	paramUUID := c.Param("uuid")

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
	uuid, err := utility.GetUUIDInternally(headers, paramUUID)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	FetchSavePostEndPoint := fmt.Sprintf(FetchUserSavedPostsEndPoint, uuid)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchSavePostEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	// If flow succeeds
	dataResponse, err := utils.PopulateDataResponseForFeed(headers, utils.GetRedisClientFromContext(c), apiCR.Response)
	if err != nil {
		utils.GenerateResponse(c, nil, false)
		return
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

func createSavePostInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	paramPostId := c.Param("post_id")
	SavePostEndPoint := fmt.Sprintf(SinglePostSaveEndPoint, paramPostId)

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

	// Add member Role in headers
	utils.AddMemberRoleToHeaders(c, response.IsCm)

	//Send Request
	utils.SendRequest(c, utils.SwarmService, SavePostEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, nil)
}
