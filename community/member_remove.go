package community

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RemoveMemberRequest struct {
	MemberIds interface{} `json:"member_ids,omitempty"`
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
		utils.GeneralAPIError(c, err.Error())
		return
	}

	user_unique_ids, err := utility.GetUsersInfo(c, removeMemberRequest.MemberIds.([]interface{}), true)

	if err != nil {
		return
	}

	removeMemberRequest.MemberIds = utils.ParseArrayToString(user_unique_ids.([]interface{}))

	// Send Request to Main service to remove member
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, RemoveMemberEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeMemberRequest)

	// Validate response and if false then return
	apiCr := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCr == nil {
		return
	}

	// If response is successfull
	var user_ids []interface{}

	// If request is for self removal, then add user id to the list
	if len(user_unique_ids.([]interface{})) == 0 {
		user_ids = append(user_unique_ids.([]interface{}), userId)
	}

	// create body for user data
	postBody := map[string]interface{}{
		"user_ids":   user_ids,
		"user_is_cm": true,
	}

	// Send request internally to delete user data
	respBytes, _, err = utils.GetRequestResponseWithoutContext(utils.SwarmService, feed.DeleteUserDataEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, postBody)

	// Validate response and log if error
	validateDeleteUserDataReponse(respBytes, err)

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

func validateDeleteUserDataReponse(respBytes []byte, err error) {

	//If API fails or any other error
	if err != nil {
		log.Println("Error while deleting user data from Swarm: ", err.Error())
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	marshal_err := api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)

	if marshal_err != nil {
		//Internal unmarshal error
		log.Println("Error while Umarshalling: ", marshal_err.Error())
	}

	if !apiCR.Success {
		//If internal api returns success as false
		log.Println("Error while deleting user data from Swarm: ", apiCR.ErrorMessage)
	}

}
