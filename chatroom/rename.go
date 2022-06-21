package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RenameChatroomRequest struct {
	ChatroomID      int64  `json:"chatroom_id" binding:"required"`
	Header          string `json:"header"`
	FirstTimeRename bool   `json:"first_time_rename"`
}

//RenameChatroom is used to rename an existing chatroom
func RenameChatroom(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/chatroom_rename POST request
	renameChatroomRequest, err := parseRenameChatroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + RenameChatroomEndPoint,
		Body:          renameChatroomRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
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

func parseRenameChatroomRequest(c *gin.Context) (*RenameChatroomRequest, error) {
	//POST body params
	var rcr RenameChatroomRequest

	if err := c.ShouldBindJSON(&rcr); err != nil {
		return nil, err
	}

	return &rcr, nil
}
