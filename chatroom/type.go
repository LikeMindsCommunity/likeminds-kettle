package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ChatroomTypeRequest struct {
	ChatroomID *int32 `json:"chatroom_id" binding:"required"`
	IsSecret   *bool  `json:"is_secret" binding:"required"`
}

//ChatroomType is used to change the type of chatroom
func ChatroomType(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Body to be sent in the api/chatroom/change_type POST request
	chatroomTypeRequest, err := parseChatroomTypeRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + ChatroomTypeEndPoint,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		Body:          chatroomTypeRequest,
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}

func parseChatroomTypeRequest(c *gin.Context) (*ChatroomTypeRequest, error) {
	//POST body params
	var ctr ChatroomTypeRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	return &ctr, nil
}
