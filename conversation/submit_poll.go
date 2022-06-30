package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PollIDObject struct {
	ID int64 `json:"id"`
}

type SubmitPollRequest struct {
	Polls          []PollIDObject `json:"polls"`
	ConversationID int64          `json:"conversation_id"`
}

//SubmitPoll is used add answer to a poll
func SubmitPoll(c *gin.Context) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/conversation/submit_poll POST request
	submitPollRequest, err := parseSubmitPollRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + SubmitPollEndPoint,
		Body:          submitPollRequest,
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

func parseSubmitPollRequest(c *gin.Context) (*SubmitPollRequest, error) {
	//POST body params
	var spr SubmitPollRequest

	if err := c.ShouldBindJSON(&spr); err != nil {
		return nil, err
	}

	return &spr, nil
}
