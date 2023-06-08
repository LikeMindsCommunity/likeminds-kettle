package search

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

// UserCreatedPostSearch is used to perform search on the user created posts
func UserCreatedPostSearch(c *gin.Context) {
	// fetch url params
	user_id := c.Param("user_id")

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, feed.VIEW_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Get user_unique_id from user_id internally
	user_id, err := utility.GetUuidInternally(utils.CreateHeaders(c, userId), user_id)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	CreatedPostSearchEndpoint := fmt.Sprintf(CreatedPostSearchEndPoint, user_id)

	//Params to be sent in the created post search api internally
	params := map[string]string{
		ParamSearch:     c.Query(ParamSearch),
		ParamSearchType: c.Query(ParamSearchType),
		ParamPage:       c.Query(ParamPage),
		ParamPageSize:   c.Query(ParamPageSize),
		ParamUserIsCm:   fmt.Sprint(response.IsCm),
	}

	//Send Request
	// utils.SendRequest(c, utils.SwarmService, CreatedPostSearchEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, CreatedPostSearchEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

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
			if user_unique_id, ok := post_data.(map[string]interface{})["user_id"]; ok {
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
