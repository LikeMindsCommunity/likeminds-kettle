package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type EditPendingPostRequest struct {
	EditPostRequest
	UUIDs []string `json:"uuids"`
}

// CreatePendingPost is used to create a new post
func CreatePendingPost(c *gin.Context) {
	PendingPost(c, utils.POSTMethod)
}

// EditPendingPost is used to create a new post
func EditPendingPost(c *gin.Context) {
	PendingPost(c, utils.PUTMethod)
}

// Pending post method handles pending post objects
func PendingPost(c *gin.Context, method int) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.POSTMethod:
		createPendingPostInternal(c, userId)

	case utils.PUTMethod:
		editPendingPostInternal(c, userId)
	}
}

// Exposed method to create a pending post for review
func createPendingPostInternal(c *gin.Context, userId string) {
	// Use Create post body params to create Pending post
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

func editPendingPostInternal(c *gin.Context, userId string) {
	pendingPostId := c.Param("pending_post_id")

	editPostEndPoint := fmt.Sprintf(EditPendingPostEndPoint, pendingPostId)

	// Body to be sent in the /post/pending/<pending_post_id> PUT request
	eppr, err := parseEditPostRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Get tagged users from text
	taggedUsers := getTaggedUsersFromText(utils.CreateHeaders(c, userId), eppr.Text)

	editPendingPostRequest := EditPendingPostRequest{
		EditPostRequest: *eppr,
		UUIDs:           taggedUsers,
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, editPostEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editPendingPostRequest)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds populate post data
	dataResponse := apiCR.Response
	dataResponse = populatePostDataResponse(c, dataResponse)

	//Generate Response
	utils.GenerateResponse(c, dataResponse, true)
}
