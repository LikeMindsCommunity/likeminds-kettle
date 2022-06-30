package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UpdateLastSeenEventRequest struct {
	ConversationID int64 `json:"conversation_id"`
}

//UpdateLastSeenEvent is used mark last seen for an event
func UpdateLastSeenEvent(c *gin.Context) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/conversation/event/update_last_seen_event POST request
	updateLastSeenEventRequest, err := parseUpdateLastSeenEventRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + UpdateLastSeenEventEndPoint,
		Body:          updateLastSeenEventRequest,
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

func parseUpdateLastSeenEventRequest(c *gin.Context) (*UpdateLastSeenEventRequest, error) {
	//POST body params
	var ulser UpdateLastSeenEventRequest

	if err := c.ShouldBindJSON(&ulser); err != nil {
		return nil, err
	}

	return &ulser, nil
}
