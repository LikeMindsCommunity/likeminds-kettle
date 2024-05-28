package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type faq struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// Request body params for adding event faq
type AddEventFAQRequest struct {
	ChatroomId string `json:"chatroom_id" binding:"required"`
	FAQ        []faq  `json:"faq,omitempty"`
}

// function to parse POST body params for AddEventFAQRequest
func parseAddEventFAQRequest(c *gin.Context) (*AddEventFAQRequest, error) {

	var ahr AddEventFAQRequest

	if err := c.ShouldBindJSON(&ahr); err != nil {
		return nil, err
	}

	return &ahr, nil
}

// Exposed function to add event faq
func AddEventFAQ(c *gin.Context) {

	userId := user.GetRequestingUserId(c)

	if userId == "" {
		return
	}

	// POST body params
	AddEventFAQRequest, err := parseAddEventFAQRequest(c)

	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to api/chatroom/event/add_or_update_event_faq
	utils.SendRequest(c, utils.CoreService, AddEventFAQEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, AddEventFAQRequest)
}
