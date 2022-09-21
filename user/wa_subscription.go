package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

//WASubscriptionRequest
type WASubscriptionRequest struct {
	Id             string `json:"id" binding:"required"`
	Created        string `json:"created" binding:"required"`
	ConversationId string `json:"conversation_id" binding:"required"`
	TicketId       string `json:"ticket_id" binding:"required"`
	Text           string `json:"text"`
	Type           string `json:"type" binding:"required"`
	Data           string `json:"data"`
	Timestamp      string `json:"timestamp" binding:"required"`
	Owner          bool   `json:"owner" binding:"required"`
	EventType      string `json:"event_type" binding:"required"`
	StatusString   string `json:"statusString"`
	AvatarURL      string `json:"avatarUrl"`
	AssignedId     string `json:"assignedId"`
	OperatorName   string `json:"operationName"`
	OperatorEmail  string `json:"operatorEmail"`
	WaId           string `json:"waId" binding:"required"`
	SenderName     string `json:"senderName" binding:"required"`
}

//WASubscription is used to update the WA subscription of a user
func WASubscription(c *gin.Context) {

	//Body to be sent in the auto follow for all members api internally
	waSubscriptionRequest, err := parseWASubscriptionRequst(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, WASubscriptionEndPoint, utils.POSTRequestRawBody, nil, nil, waSubscriptionRequest)
}

func parseWASubscriptionRequst(c *gin.Context) (*WASubscriptionRequest, error) {
	//POST body params
	var wasr WASubscriptionRequest

	if err := c.ShouldBindJSON(&wasr); err != nil {
		return nil, err
	}

	return &wasr, nil
}
