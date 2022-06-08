package chatroom

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//AutoFollowMembersRequest
type AutoFollowMembersRequest struct {
	ChatroomId          int32 `json:"chatroom_id" binding:"required"`
	AutoFollowDone      bool  `json:"auto_follow_done" binding:"required"`
	IncludeMembersLater bool  `json:"include_members_later" binding:"required"`
}

//AutoFollowMembers is used to enable auto follow members for a chatroom
func AutoFollowMembers(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Body to be sent in the api/chatroom/auto_follow_for_all_members POST request
	autoFollowMembersRequest, err := parseAutoFollowMembersRequst(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + AutoFollowMembersEndPoint,
		Body:          autoFollowMembersRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
	}

	if !apiCR.Success {
		//If api/chatroom/auto_follow_for_all_members returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response with api/chatroom/auto_follow_for_all_members response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}

func parseAutoFollowMembersRequst(c *gin.Context) (*AutoFollowMembersRequest, error) {
	//POST body params
	var afmr AutoFollowMembersRequest

	if err := c.ShouldBindJSON(&afmr); err != nil {
		return nil, err
	}

	return &afmr, nil
}
