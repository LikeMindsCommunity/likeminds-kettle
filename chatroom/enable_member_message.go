package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//EnableMemberMessageRequest | member message setting schema
type EnableMemberMessageRequest struct {
	ChatroomId int32 `json:"chatroom_id" binding:"required"`
	Value      bool  `json:"value" binding:"required"`
}

//EnableMemberMessage is used to enable member message settings in chatroom
func EnableMemberMessage(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Body to be sent in the api/chatroom/enable_member_message POST request
	enableMemberMessageRequest, err := parseEnableMemberMessageRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + EnableMemberMessageEndPoint,
		Body:          enableMemberMessageRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}

func parseEnableMemberMessageRequest(c *gin.Context) (*EnableMemberMessageRequest, error) {
	//POST body params
	var emmr EnableMemberMessageRequest

	if err := c.ShouldBindJSON(&emmr); err != nil {
		return nil, err
	}

	return &emmr, nil
}
