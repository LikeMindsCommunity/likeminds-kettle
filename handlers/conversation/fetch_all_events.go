package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// FetchAllEvents is used to fetch all events
func FetchAllEvents(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch all events api internally
	params := map[string]string{
		ParamPastEvents:      c.Query(ParamPastEvents),
		ParamAttendingStatus: c.Query(ParamAttendingStatus),
		ParamPage:            c.Query(ParamPage),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EventFetchAllEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
