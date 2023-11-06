package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

// FetchUserCreatedPost is used to fetch posts created by a user
func FetchUserCreatedPosts(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the /user/<user_id>/post request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Access query params and url generation
	user_id := c.Param("user_id")

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
	user_id, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), user_id)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	UserCreatedPostsEndPoint := fmt.Sprintf(FetchUserCreatedPostsEndPoint, user_id)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, UserCreatedPostsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["posts"]; ok {
		posts := value.([]interface{})
		user_ids := []string{}

		//Fetch posts user id
		for _, post_data := range posts {
			if user_unique_id, ok := post_data.(map[string]interface{})["uuid"]; ok {
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
	utils.GenerateResponse(c, dataResponse, false)
}
