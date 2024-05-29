package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type Instructors struct {
	About string `json:"about"`
	Url   string `json:"url"`
	Name  string `json:"name"`
}

// Request body params for adding event instructors
type AddEventInstructorsRequest struct {
	ChatroomId  string        `json:"chatroom_id" binding:"required"`
	Instructors []Instructors `json:"instructors,omitempty"`
}

// function to parse POST body params for AddEventInstructorsRequest
func parseAddEventInstructorsRequest(c *gin.Context) (*AddEventInstructorsRequest, error) {

	var aer AddEventInstructorsRequest

	if err := c.ShouldBindJSON(&aer); err != nil {
		return nil, err
	}

	return &aer, nil
}

// Exposed function to add event instructors
func AddEventInstructors(c *gin.Context) {

	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// POST body params
	AddEventInstructorsRequest, err := parseAddEventInstructorsRequest(c)

	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to api/chatroom/event/add_or_update_instructor
	utils.SendRequest(c, utils.CoreService, AddEventInstructorsEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, AddEventInstructorsRequest)
}
