package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UpdateLastSeenEventRequest struct {
	ConversationID int64 `json:"conversation_id"`
}

//UpdateLastSeenEvent is used mark last seen for an event
func UpdateLastSeenEvent(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the update last seen event api internally
	updateLastSeenEventRequest, err := parseUpdateLastSeenEventRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, UpdateLastSeenEventEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, updateLastSeenEventRequest)
}

func parseUpdateLastSeenEventRequest(c *gin.Context) (*UpdateLastSeenEventRequest, error) {
	//POST body params
	var ulser UpdateLastSeenEventRequest

	if err := c.ShouldBindJSON(&ulser); err != nil {
		return nil, err
	}

	return &ulser, nil
}
