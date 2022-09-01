package community

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type RemoveMemberRequest struct {
	UserId int64 `json:"user_id" binding:"required"`
	IsCM   bool  `json:"is_cm"`
}

type InternalRemoveMemberRequest struct {
	MemberIds []string `json:"member_ids"`
}

//RemoveMember is used to remove a member of CM from community
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

	is_cm := removeMemberRequest.IsCM

	if !is_cm {
		//If is_cm is missing or false, call remove from member api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, RemoveMemberEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, updateRemoveMemberRequest(removeMemberRequest))
	} else {
		//else, call remove community manager api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, RemoveCMEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeMemberRequest)
	}
}

func parseRemoveMemberRequest(c *gin.Context) (*RemoveMemberRequest, error) {
	//POST body params
	var rmr RemoveMemberRequest

	if err := c.ShouldBindJSON(&rmr); err != nil {
		return nil, err
	}

	return &rmr, nil
}

func updateRemoveMemberRequest(rmr *RemoveMemberRequest) interface{} {
	var updatedRmr InternalRemoveMemberRequest

	updatedRmr.MemberIds = []string{strconv.Itoa(int(rmr.UserId))}

	return updatedRmr
}
