package community

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RemoveMemberRequest struct {
	MemberIds []string `json:"member_ids,omitempty"`
	TagID     int32    `json:"tag_id"`
	Reason    string   `json:"reason"`
}

// RemoveMember is used to remove a member from community
func RemoveMember(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the remove member request
	removeMemberRequest, err := parseRemoveMemberRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request to Main service to remove member
	respbytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, RemoveMemberEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeMemberRequest)

	// Validate response and if false then return
	apiCr := utils.ValidateClientResponse(c, respbytes, statusCode)
	if apiCr == nil {
		return
	}

	// If response is successfull
	user_ids := removeMemberRequest.MemberIds

	// If request is for self removal, then add user id to the list
	if len(user_ids) == 0 {
		user_ids = append(user_ids, userId)
	}

	// create body for user data
	postBody := map[string]interface{}{
		"user_ids":   user_ids,
		"is_user_cm": true,
	}

	// Send request internally to delete user data
	_, _, err = utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.DeleteUserDataEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, postBody)
	if err != nil {
		log.Println("Error while deleting user data from Swarm: ", err)
	}

	//Generate response
	utils.GenerateResponse(c, apiCr.Response)
}

func parseRemoveMemberRequest(c *gin.Context) (*RemoveMemberRequest, error) {
	//POST body params
	var rmr RemoveMemberRequest

	if err := c.ShouldBindJSON(&rmr); err != nil {
		return nil, err
	}

	return &rmr, nil
}
