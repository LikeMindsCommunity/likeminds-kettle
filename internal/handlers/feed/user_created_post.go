package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// FetchUserCreatedPost is used to fetch posts created by a user
func FetchUserCreatedPosts(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent in the /user/<user_id>/post request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Access query params and url generation
	paramUserId := c.Param("user_id")

	//Fetch member access to view post
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
	user_id, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), paramUserId)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	UserCreatedPostsEndPoint := fmt.Sprintf(FetchUserCreatedPostsEndPoint, user_id)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, UserCreatedPostsEndPoint, utils.GETRequest, headers, params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse, err := utils.PopulateDataResponseForFeed(headers, utils.GetRedisClientFromContext(c), apiCR.Response)
	if err != nil {
		utils.GenerateResponse(c, nil, false)
		return
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
