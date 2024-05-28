package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type PollIDObject struct {
	ID interface{} `json:"id"`
}

type SubmitPollRequest struct {
	Polls          []PollIDObject `json:"polls"`
	ConversationID interface{}    `json:"conversation_id"`
}

// SubmitPoll is used add answer to a poll
func SubmitPoll(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the submit poll api internally
	submitPollRequest, err := parseSubmitPollRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, SubmitPollEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, submitPollRequest)
}

func parseSubmitPollRequest(c *gin.Context) (*SubmitPollRequest, error) {
	//POST body params
	var spr SubmitPollRequest

	if err := c.ShouldBindJSON(&spr); err != nil {
		return nil, err
	}

	if spr.ConversationID == nil {
		spr.ConversationID = utils.ParseInterfaceToString(spr.ConversationID)
	}

	return &spr, nil
}
