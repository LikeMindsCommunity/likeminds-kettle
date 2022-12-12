package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateUserActivityRequest struct {
	Action string `json:"action" binding:"required"`
}

func parseCreateUserActivityRequest(c *gin.Context) (*CreateUserActivityRequest, error) {
	//POST body params
	var cuar CreateUserActivityRequest

	if err := c.ShouldBindJSON(&cuar); err != nil {
		return nil, err
	}

	return &cuar, nil
}

// CreateUserActivity is used to create activity for a user
func CreateUserActivity(c *gin.Context) {
	UserActivity(c, utils.PUTMethod)
}

// UserActivity method handles post like objects
func UserActivity(c *gin.Context, method int) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	user_id := c.Param("user_id")
	UserActivityEndPoint := fmt.Sprintf(SingleUserActivityEndPoint, user_id)

	//Send request
	switch method {
	case utils.PUTMethod:
		createUserActivityInternal(c, userId, UserActivityEndPoint)

	}
}

func createUserActivityInternal(c *gin.Context, userId string, EndPoint string) {
	//Body to be sent in the /user/<user_id>/activity POST request
	createUserActivityRequest, err := parseCreateUserActivityRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
	}

	//Fetch member access to create post
	success, response := user.FetchMemberAccess(c, CREATE_ACTIVITY_ACTION)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, EndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createUserActivityRequest)
}
