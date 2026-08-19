package conversation

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type EventAttendedRequest struct {
	ConversationID interface{} `json:"conversation_id"`
}

// EventAttended is used to send attendence of a user
func EventAttended(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the event attended api internally
	eventAttendedRequest, err := parseEventAttendedRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EventAttendedEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, eventAttendedRequest)
}

func parseEventAttendedRequest(c *gin.Context) (*EventAttendedRequest, error) {
	//POST body params
	var ear EventAttendedRequest

	if err := c.ShouldBindJSON(&ear); err != nil {
		return nil, err
	}

	if ear.ConversationID != nil {
		ear.ConversationID = utils.ParseInterfaceToString(ear.ConversationID)
	}

	return &ear, nil
}
