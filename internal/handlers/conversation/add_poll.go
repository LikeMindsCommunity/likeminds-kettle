package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type AddPollRequest struct {
	Poll           PollObject  `json:"poll"`
	ConversationID interface{} `json:"conversation_id"`
}

// AddPoll is used add answer to a poll
func AddPoll(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the add poll api internally
	addPollRequest, err := parseAddPollRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, AddPollEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, addPollRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}

func parseAddPollRequest(c *gin.Context) (*AddPollRequest, error) {
	//POST body params
	var apr AddPollRequest

	if err := c.ShouldBindJSON(&apr); err != nil {
		return nil, err
	}

	if apr.ConversationID == nil {
		apr.ConversationID = utils.ParseInterfaceToString(apr.ConversationID)
	}

	return &apr, nil
}
