package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type SetTopicRequest struct {
	ChatroomID     interface{} `json:"chatroom_id"`
	ConversationID interface{} `json:"conversation_id"`
}

// SetTopic is used to set topic for conversation
func SetTopic(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the set topic api internally
	setTopicRequest, err := parseSetTopicRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, SetTopicEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, setTopicRequest)
}

func parseSetTopicRequest(c *gin.Context) (*SetTopicRequest, error) {
	//POST body params
	var str SetTopicRequest

	if err := c.ShouldBindJSON(&str); err != nil {
		return nil, err
	}

	if str.ChatroomID == nil {
		str.ChatroomID = utils.ParseInterfaceToString(str.ChatroomID)
	}

	if str.ConversationID == nil {
		str.ConversationID = utils.ParseInterfaceToString(str.ConversationID)
	}

	return &str, nil
}
