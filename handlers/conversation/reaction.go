package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddReactionRequest struct {
	ChatroomID     interface{} `json:"chatroom_id,omitempty"`
	Reaction       string      `json:"reaction"`
	ConversationID interface{} `json:"conversation_id,omitempty"`
}

type RemoveReactionRequest struct {
	ChatroomID     interface{} `json:"chatroom_id,omitempty"`
	ConversationID interface{} `json:"conversation_id,omitempty"`
}

// AddReaction is used to add reaction to specific conversation
func AddReaction(c *gin.Context) {
	Reaction(c, utils.PUTMethod)
}

// RemoveReaction is used to delete reaction from a specific conversation
func RemoveReaction(c *gin.Context) {
	Reaction(c, utils.DELETEMethod)
}

// Reaction method handles reaction on a conversation object
func Reaction(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.PUTMethod:

		addReactionInternal(c, userId)

	case utils.DELETEMethod:

		removeReactionInternal(c, userId)
	}
}

func parseAddReactionRequest(c *gin.Context) (*AddReactionRequest, error) {
	//POST body params
	var arr AddReactionRequest

	if err := c.ShouldBindJSON(&arr); err != nil {
		return nil, err
	}

	// Parse ConversationId and ChatroomId to string if not null
	if arr.ConversationID != nil {
		arr.ConversationID = utils.ParseInterfaceToString(arr.ConversationID)
	}

	if arr.ChatroomID != nil {
		arr.ChatroomID = utils.ParseInterfaceToString(arr.ChatroomID)
	}

	return &arr, nil
}

func parseRemoveReactionRequest(c *gin.Context) (*RemoveReactionRequest, error) {
	//POST body params
	var rrr RemoveReactionRequest

	if err := c.ShouldBindJSON(&rrr); err != nil {
		return nil, err
	}

	if rrr.ConversationID != nil {
		rrr.ConversationID = utils.ParseInterfaceToString(rrr.ConversationID)
	}

	if rrr.ChatroomID != nil {
		rrr.ChatroomID = utils.ParseInterfaceToString(rrr.ChatroomID)
	}

	return &rrr, nil
}

func addReactionInternal(c *gin.Context, userId string) {

	//Body to be sent in the add reaction api internally
	addReactionRequest, err := parseAddReactionRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, AddReactionEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, addReactionRequest)
}

func removeReactionInternal(c *gin.Context, userId string) {

	//Body to be sent in the remove reaction api internally
	removeReactionRequest, err := parseRemoveReactionRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, RemoveReactionEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeReactionRequest)
}
