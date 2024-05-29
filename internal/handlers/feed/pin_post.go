package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// PinPost is used to pin a post
func PinPost(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	post_id := c.Param("post_id")
	PinPostEndPoint := fmt.Sprintf(SinglePostPinEndPoint, post_id)

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, PIN_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, PinPostEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, nil)
}
