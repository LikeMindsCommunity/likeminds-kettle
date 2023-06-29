package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type testimonials struct {
	MemberName  string `json:"member_name"`
	Testimonial string `json:"testimonial"`
	Url         string `json:"url"`
}

// Request body params for adding event testimonials
type addEventTestimonialsRequest struct {
	ChatroomId   string         `json:"chatroom_id"`
	Testimonials []testimonials `json:"testimonials"`
}

// function to parse POST body params for addEventTestimonialsRequest
func parseAddEventTestimonialsRequest(c *gin.Context) (*addEventTestimonialsRequest, error) {

	var aer addEventTestimonialsRequest

	if err := c.ShouldBindJSON(&aer); err != nil {
		return nil, err
	}

	return &aer, nil
}

// Exposed function to add event testimonials
func AddEventTestimonials(c *gin.Context) {

	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// POST body params
	addEventTestimonialsRequest, err := parseAddEventTestimonialsRequest(c)

	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to api/chatroom/event/add_or_update_member_testimonial
	utils.SendRequest(c, utils.CoreService, AddEventTestimonialsEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, addEventTestimonialsRequest)
}
