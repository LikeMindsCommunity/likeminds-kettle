package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type MuteChatroomRequest struct {
	ChatroomID int64 `json:"chatroom_id" binding:"required"`
	Value      bool  `json:"value" binding:"required"`
}

//MuteChatroom is used to mute a specifid chatroom
func MuteChatroom(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/chatroom_mute POST request
	muteChatroomRequest, err := parseMuteChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + MuteChatroomEndPoint,
		Body:          muteChatroomRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
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

func parseMuteChatroomRequest(c *gin.Context) (*MuteChatroomRequest, error) {
	//POST body params
	var mcr MuteChatroomRequest

	if err := c.ShouldBindJSON(&mcr); err != nil {
		return nil, err
	}

	return &mcr, nil
}
