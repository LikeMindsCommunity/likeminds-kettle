package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type EventAttendedRequest struct {
	ConversationID int64 `json:"conversation_id"`
}

//EventAttended is used to send attendence of a user
func EventAttended(c *gin.Context) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/conversation/event/attended POST request
	eventAttendedRequest, err := parseEventAttendedRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + EventAttendedEndPoint,
		Body:          eventAttendedRequest,
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

func parseEventAttendedRequest(c *gin.Context) (*EventAttendedRequest, error) {
	//POST body params
	var ear EventAttendedRequest

	if err := c.ShouldBindJSON(&ear); err != nil {
		return nil, err
	}

	return &ear, nil
}
