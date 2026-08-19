package feed

import (
	"fmt"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/utility"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func FetchUserProfileActivity(c *gin.Context) {

	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	paramUUID := c.Param("uuid")

	//Get user_unique_id from uuid internally
	uuid, err := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), paramUUID)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

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

	//Url generation
	UserActivityEndPoint := fmt.Sprintf(constants.UserProfileActivityEndPoint, uuid)

	//Send Request
	utils.SendRequest(c, utils.SwarmService, UserActivityEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
