package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PinChatroomRequest struct {
	ChatroomID int64 `json:"chatroom_id" binding:"required"`
	Value      *bool `json:"value" binding:"required"`
	Notify     bool  `json:"notify"`
}

//PinChatroom is used to create a pin a chatroom
func PinChatroom(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Body to be sent in the api/chatroom/pin POST request
	pinChatroomRequest, err := parsePinChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + PinChatroomEndPoint,
		Body:          pinChatroomRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse Response
	utils.ParseResponse(c, respBytes)
}

func parsePinChatroomRequest(c *gin.Context) (*PinChatroomRequest, error) {
	//POST body params
	var pcr PinChatroomRequest

	if err := c.ShouldBindJSON(&pcr); err != nil {
		return nil, err
	}

	return &pcr, nil
}
