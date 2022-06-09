package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type DeleteChatroomRequest struct {
	ChatroomID int64  `json:"chatroom_id"`
	TagID      int32  `json:"tag_id"`
	Reason     string `json:"reason"`
}

//DeleteChatroom is used to delete an existing chatroom
func DeleteChatroom(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Body to be sent in the api/chatroom_delete POST request
	deleteChatroomRequest, err := parseDeleteChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + DeleteChatroomEndPoint,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		Body:          deleteChatroomRequest,
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeFormUrlEncoded)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}

func parseDeleteChatroomRequest(c *gin.Context) (*DeleteChatroomRequest, error) {
	//POST body params
	var dcr DeleteChatroomRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}
