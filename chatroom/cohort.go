package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddCohortChatroom struct {
	ChatroomID         interface{} `json:"chatroom_id" binding:"required"`
	CohortIDs          []int       `json:"cohort_ids"`
	AddExistingMembers bool        `json:"add_existing_members"`
}

type RemoveCohortChatroom struct {
	ChatroomID interface{} `json:"chatroom_id" binding:"required"`
	CohortID   int         `json:"cohort_id"`
}

// Add cohort to chatroom
func AddCohortToChatroom(c *gin.Context) {
	ChatroomCohort(c, utils.POSTMethod)
}

// Remove cohort from chatroom
func RemoveCohortFromChatroom(c *gin.Context) {
	ChatroomCohort(c, utils.DELETEMethod)
}

func ChatroomCohort(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Send request
	switch method {
	case utils.POSTMethod:

		// Body to be sent in the add cohort POST request
		addCohortRequest, err := parseAddCohortRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, AddCohortsToChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, addCohortRequest)

	case utils.DELETEMethod:
		// Body to be sent in the add cohort POST request
		removeCohortRequest, err := parseRemoveCohortRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, RemoveCohortFromChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, removeCohortRequest)
	}
}

func parseAddCohortRequest(c *gin.Context) (*AddCohortChatroom, error) {
	// POST body params
	var acc AddCohortChatroom

	if err := c.ShouldBindJSON(&acc); err != nil {
		return nil, err
	}

	if acc.ChatroomID != nil {
		acc.ChatroomID = utils.ParseInterfaceToString(acc.ChatroomID)
	}

	return &acc, nil
}

func parseRemoveCohortRequest(c *gin.Context) (*RemoveCohortChatroom, error) {
	// POST body params
	var rcc RemoveCohortChatroom

	if err := c.ShouldBindJSON(&rcc); err != nil {
		return nil, err
	}

	if rcc.ChatroomID != nil {
		rcc.ChatroomID = utils.ParseInterfaceToString(rcc.ChatroomID)
	}

	return &rcc, nil
}
