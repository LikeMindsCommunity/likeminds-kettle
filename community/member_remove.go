package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RemoveMemberRequest struct {
	MemberIds []string `json:"member_ids,omitempty"`
	TagID     int32    `json:"tag_id"`
	Reason    string   `json:"reason"`
}

//RemoveMember is used to remove a member from community
func RemoveMember(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the remove member request
	removeMemberRequest, err := parseRemoveMemberRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, RemoveMemberEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeMemberRequest)
}

func parseRemoveMemberRequest(c *gin.Context) (*RemoveMemberRequest, error) {
	//POST body params
	var rmr RemoveMemberRequest

	if err := c.ShouldBindJSON(&rmr); err != nil {
		return nil, err
	}

	return &rmr, nil
}
