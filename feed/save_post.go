package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
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
	FetchSavePostEndPoint := fmt.Sprintf(FetchUserSavedPostsEndPoint, user_id)

	//Params to be sent in the /user/<user_id>/save request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Fetch member access to view post likes
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION)
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
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchSavePostEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

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
		success, user_data := user.FetchMemberMeta(c, user_ids)
		if !success {
			return
		}

		//Update user data for posts
		for post_index, post_data := range posts {
			if user_unique_id, ok := post_data.(map[string]interface{})["user_id"]; ok {
				for _, member := range user_data.Members {
					if member.UserUniqueId == user_unique_id {
						dataResponse["posts"].([]interface{})[post_index].(map[string]interface{})["user"] = member
					}
				}
				delete(dataResponse["posts"].([]interface{})[post_index].(map[string]interface{}), "user_id")
			}
		}
	}

	//Send response
	utils.GenerateResponse(c, dataResponse)
}

func createSavePostInternal(c *gin.Context, userId string) {
	//Access query params and url generation
	post_id := c.Param("post_id")
	SavePostEndPoint := fmt.Sprintf(SinglePostSaveEndPoint, post_id)

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, SAVE_POST_ACTION)
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
