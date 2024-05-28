package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type EventAttendRequest struct {
	ConversationID  interface{} `json:"conversation_id"`
	AttendingStatus bool        `json:"attending_status"`
}

// EventAttend is used mark event as attend
func EventAttend(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the attend event api internally
	eventAttendRequest, err := parseEventAttendRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EventAttendEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, eventAttendRequest)
}

func parseEventAttendRequest(c *gin.Context) (*EventAttendRequest, error) {
	//POST body params
	var ear EventAttendRequest

	if err := c.ShouldBindJSON(&ear); err != nil {
		return nil, err
	}

	if ear.ConversationID != nil {
		ear.ConversationID = utils.ParseInterfaceToString(ear.ConversationID)
	}

	return &ear, nil
}
