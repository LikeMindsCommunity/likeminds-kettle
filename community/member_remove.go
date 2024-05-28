package community

import (
	"fmt"

	log "github.com/nateshr/likeminds-authentication/logging"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/handlers/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RemoveMemberRequest struct {
	MemberIds interface{} `json:"member_ids,omitempty"`
	UUIDs     interface{} `json:"uuids,omitempty"`
	TagID     int32       `json:"tag_id"`
	Reason    string      `json:"reason"`
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
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	if removeMemberRequest.MemberIds != nil {
		user_unique_ids, err := utility.FetchUserUniqueIdsFromAnyUserIds(utils.CreateHeaders(c, userId), removeMemberRequest.MemberIds)
		if err != nil {
			utils.GeneralAPIError(c, fmt.Sprintf("Error while fetching user info: %s", err.Error()))
			return
		}

		if len(user_unique_ids) == 0 {
			utils.GeneralBadRequestError(c, utils.ErrorNoUserFoundWithGivenIds)
			return
		}

		removeMemberRequest.MemberIds = user_unique_ids
	}

	// Send Request to caravan service to remove member
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, RemoveMemberEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeMemberRequest)

	// Validate response and if false then return
	apiCr := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCr == nil {
		return
	}

	// If request is for self removal, then add user id to the list
	if removeMemberRequest.MemberIds == nil {
		removeMemberRequest.MemberIds = []string{userId}
	}
	// create body for user data
	postBody := map[string]interface{}{
		"user_ids":   removeMemberRequest.MemberIds,
		"user_is_cm": true,
	}

	// Send request internally to delete user data
	respBytes, _, err = utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.DeleteUserDataEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, postBody)

	// Validate response and log if error
	validateDeleteUserDataReponse(respBytes, err)

	//Generate response
	utils.GenerateResponse(c, apiCr.Response, false)
}

func parseRemoveMemberRequest(c *gin.Context) (*RemoveMemberRequest, error) {
	//POST body params
	var rmr RemoveMemberRequest

	if err := c.ShouldBindJSON(&rmr); err != nil {
		return nil, err
	}

	if rmr.UUIDs != nil {
		rmr.MemberIds = rmr.UUIDs
	}

	return &rmr, nil
}

func validateDeleteUserDataReponse(respBytes []byte, err error) {

	//If API fails or any other error
	if err != nil {
		log.Error(fmt.Sprintf("Error while deleting user data from Swarm: %s", err.Error()))
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	marshal_err := api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)

	if marshal_err != nil {
		//Internal unmarshal error
		log.Error(fmt.Sprintf("Error while Umarshalling: %s", marshal_err.Error()))
	}

	if !apiCR.Success {
		//If internal api returns success as false
		log.Error(fmt.Sprintf("Error while deleting user data from Swarm: %s", apiCR.ErrorMessage))
	}

}
