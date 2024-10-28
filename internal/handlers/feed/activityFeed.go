package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func FetchUserProfileActivity(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	user_id := c.Param("user_id")

	//Get user_unique_id from user_id internally
	user_id, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), user_id)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Url generation
	UserActivityEndPoint := fmt.Sprintf(SingleUserActivityEndPoint, user_id)

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, VIEW_USER_ACTIVITY, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	utils.AddMemberRoleToHeaders(c, response.IsCm)

	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, UserActivityEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
