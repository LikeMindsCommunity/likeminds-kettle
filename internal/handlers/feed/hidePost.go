package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// HidePost is used to Hide a post
func HidePost(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	post_id := c.Param("post_id")
	HidePostEndpoint := fmt.Sprintf(SinglePostHideEndPoint, post_id)

	// Fetch member access
	success, response := user.FetchMemberAccess(c, HIDE_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	// add CM role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, HidePostEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, nil)
}
