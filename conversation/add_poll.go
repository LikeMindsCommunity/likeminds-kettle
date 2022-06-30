package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddPollRequest struct {
	Poll           PollObject `json:"poll"`
	ConversationID int64      `json:"conversation_id"`
}

//AddPoll is used add answer to a poll
func AddPoll(c *gin.Context) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Body to be sent in the api/conversation/add_poll POST request
	addPollRequest, err := parseAddPollRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + AddPollEndPoint,
		Body:          addPollRequest,
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

func parseAddPollRequest(c *gin.Context) (*AddPollRequest, error) {
	//POST body params
	var apr AddPollRequest

	if err := c.ShouldBindJSON(&apr); err != nil {
		return nil, err
	}

	return &apr, nil
}
