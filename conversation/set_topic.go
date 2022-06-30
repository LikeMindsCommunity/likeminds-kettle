package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type SetTopicRequest struct {
	ChatroomID     int64 `json:"chatroom_id"`
	ConversationID int64 `json:"conversation_id"`
}

//SetTopic is used to set topic for conversation
func SetTopic(c *gin.Context) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/conversation/set_topic POST request
	setTopicRequest, err := parseSetTopicRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + SetTopicEndPoint,
		Body:          setTopicRequest,
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

func parseSetTopicRequest(c *gin.Context) (*SetTopicRequest, error) {
	//POST body params
	var str SetTopicRequest

	if err := c.ShouldBindJSON(&str); err != nil {
		return nil, err
	}

	return &str, nil
}
