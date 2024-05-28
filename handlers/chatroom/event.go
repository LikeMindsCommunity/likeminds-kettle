package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// Request body params to create an event
type CreateEventRequest struct {
	Type                   int32   `json:"type"`
	Header                 string  `json:"header"`
	Title                  string  `json:"title"`
	About                  string  `json:"about"`
	OnlineLinkType         int32   `json:"online_link_type"`
	OnlineLink             string  `json:"online_link"`
	OnlineLinkId           string  `json:"online_link_id"`
	OnlineLinkPassword     string  `json:"online_link_password"`
	Location               string  `json:"location,omitempty"`
	LocationLat            float32 `json:"location_lat,omitempty"`
	LocationLong           float32 `json:"location_long,omitempty"`
	DateTime               int64   `json:"date_time"`
	EndDate                int64   `json:"end_date,omitempty"`
	IsPaid                 bool    `json:"is_paid"`
	Access                 int32   `json:"access"`
	CoHost                 []int64 `json:"co_host,omitempty"`
	AttachmentCount        int64   `json:"attachment_count"`
	OnlineLinkEnableBefore int64   `json:"online_link_enable_before"`
	EventPaymentLink       string  `json:"event_payment_link"`
	EventWebPage           string  `json:"event_web_page"`
	WebFlowItemId          string  `json:"webflow_item_id"`
	AttendingStatus        bool    `json:"attending_status"`
	EventKind              string  `json:"event_kind"`
}

// Request body params to edit an event
type EditEventRequest struct {
	CreateEventRequest
	EventType  string  `json:"event_type" binding:"required"`
	ChatroomId string  `json:"chatroom_id" binding:"required"`
	CohortIds  []int64 `json:"cohort_ids,omitempty"`
}

// function to parse POST body params for CreateEventRequest
func parseCreateEventRequest(c *gin.Context) (*CreateEventRequest, error) {
	// POST body params
	var cer CreateEventRequest

	if err := c.ShouldBindJSON(&cer); err != nil {
		return nil, err
	}

	return &cer, nil
}

// function to parse POST body params for EditEventRequest
func parseEditEventRequest(c *gin.Context) (*EditEventRequest, error) {

	var eer EditEventRequest

	if err := c.ShouldBindJSON(&eer); err != nil {
		return nil, err
	}

	return &eer, nil
}

// FetchEvents is used to fetch all the events
func FetchEvents(c *gin.Context) {
	Event(c, utils.GETMethod)
}

// CreateEvent is used to create a new event
func CreateEvent(c *gin.Context) {
	Event(c, utils.POSTMethod)
}

// EditEvent is used to edit event details
func EditEvent(c *gin.Context) {
	Event(c, utils.PUTMethod)
}

// Event method handles event objects
func Event(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.GETMethod:

		getEventsInternal(c, userId)

	case utils.POSTMethod:

		createEventInternal(c, userId)

	case utils.PUTMethod:

		editEventInternal(c, userId)
	}
}

func getEventsInternal(c *gin.Context, userId string) {

	params := map[string]string{
		ParamPage:            c.Query(ParamPage),
		ParamPastEvents:      c.Query(ParamPastEvents),
		ParamAttendingStatus: c.Query(ParamAttendingStatus),
		ParamHasContent:      c.Query(ParamHasContent),
	}

	// Send request to /api/chatroom/event/fetch_all
	utils.SendRequest(c, utils.CoreService, FetchEventsEndPoint, utils.GETMethod, utils.CreateHeaders(c, userId), params, nil)
}

func createEventInternal(c *gin.Context, userId string) {

	// body to be sent in POST request
	createEventRequest, err := parseCreateEventRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to /api/chatroom/event/create
	utils.SendRequest(c, utils.CoreService, CreateEventEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, createEventRequest)
}

func editEventInternal(c *gin.Context, userId string) {

	// body to be sent in PUT request
	editEventRequest, err := parseEditEventRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// check for event_type and send request accordingly
	switch editEventRequest.EventType {

	// If event_type == "EDIT_EVENT"
	case EditEventType:

		// Send request to /api/chatroom/event/update
		utils.SendRequest(c, utils.CoreService, UpdateEventEndpoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, editEventRequest)
		return

	// If event_type == "UPDATE_LAST_SEEN"
	case UpdateLastSeenType:

		// Send request to /api/chatroom/event/update_last_seen
		utils.SendRequest(c, utils.CoreService, EditLastSeenEventEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, editEventRequest)
		return

	// If event_type == "EVENT_ATTEND"
	case EventAttendType:

		// Send request to /api/chatroom/event/attend
		utils.SendRequest(c, utils.CoreService, AttendEventEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, editEventRequest)
		return

	// If event_type == "EVENT_ATTENDED"
	case EventAttendedType:

		// Send request to /api/chatroom/event/attended
		utils.SendRequest(c, utils.CoreService, AttendedEventEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, editEventRequest)
		return
	}

	// If event_type is not valid, return error
	utils.GeneralBadRequestError(c, "Invalid event_type")
}
