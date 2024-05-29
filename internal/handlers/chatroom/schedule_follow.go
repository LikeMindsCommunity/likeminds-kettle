package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type ScheduleFollowRequest struct {
	ChatroomID interface{} `json:"chatroom_id"`
}

// ScheduleFollow is used to schedule follow request for particular user
func ScheduleFollow(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the schedule follow request internally
	scheduleFollowRequest, err := parseScheduleFollowRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ScheduleFollowEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, scheduleFollowRequest)
}

func parseScheduleFollowRequest(c *gin.Context) (*ScheduleFollowRequest, error) {
	//POST body params
	var sfr ScheduleFollowRequest
	if err := c.ShouldBindJSON(&sfr); err != nil {
		return nil, err
	}

	if sfr.ChatroomID != nil {
		sfr.ChatroomID = utils.ParseInterfaceToString(sfr.ChatroomID)
	}

	return &sfr, nil
}
