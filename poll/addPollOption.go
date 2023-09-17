package poll

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddPollOptionRequest struct {
	Text string `json:"text"`
}

func parseAddPollOptionRequest(c *gin.Context) (*AddPollOptionRequest, error) {
	//POST body params
	var apor AddPollOptionRequest

	if err := c.ShouldBindJSON(&apor); err != nil {
		return nil, err
	}

	return &apor, nil
}

// AddPollOption is used to add an option on poll
func AddPollOption(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the /poll/poll_id PUT request
	addPollOptionRequest, err := parseAddPollOptionRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Access query params and url generation
	pollId := c.Param("poll_id")
	AddPollOptionEndPoint := fmt.Sprintf(PollEndPoint, pollId)

	//Fetch member access to add poll option
	success, response := user.FetchMemberAccess(c, feed.IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, AddPollOptionEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, addPollOptionRequest)
}
