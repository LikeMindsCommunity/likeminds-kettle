package feed

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Exposed method to create a pending post for review
func CreatePendingPost(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the /post POST request
	cppr, err := parseCreatePostRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, CREATE_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	// Send request to "/post/pending" of swarm service
	utils.SendRequest(c, utils.SwarmService, CreatePendingPostEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, cppr)
}
