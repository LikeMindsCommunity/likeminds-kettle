package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, AddPollEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, addPollRequest)
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
