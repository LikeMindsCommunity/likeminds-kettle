package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type highlights struct {
	Highlight string `json:"highlight"`
	Url       string `json:"url"`
}

// Request body params for adding event highlights
type AddEventHighlightsRequest struct {
	ChatroomId string       `json:"chatroom_id" binding:"required"`
	Highlights []highlights `json:"highlights,omitempty"`
}

// Function to parse POST body params for AddEventHighlightsRequest
func parseAddEventHighlightsRequest(c *gin.Context) (*AddEventHighlightsRequest, error) {

	var ahr AddEventHighlightsRequest

	if err := c.ShouldBindJSON(&ahr); err != nil {
		return nil, err
	}

	return &ahr, nil
}

// Exposed function to add event highlights
func AddEventHighlights(c *gin.Context) {

	userId := user.GetRequestingUserId(c)

	if userId == "" {
		return
	}

	// POST body params
	AddEventHighlightsRequest, err := parseAddEventHighlightsRequest(c)

	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to api/chatroom/event/add_or_update_highlight
	utils.SendRequest(c, utils.CoreService, AddEventHighlightsEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, AddEventHighlightsRequest)
}
