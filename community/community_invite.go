package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CommunityInviteRequest struct {
	InviteType string `json:"type" binding:"required"`
	EmailID    string `json:"email_id"`
	MobileNo   string `json:"mobile_no"`
	Text       string `json:"text"`
	LinkType   string `json:"link_type"`
}

func SendCommunityInvite(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)

	if userId == "" {
		return
	}

	// Body to be sent in the api/community/invite POST request
	communityInviteRequest, err := parseCommunityInviteRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, SendCommunityInviteEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, communityInviteRequest)

}

func parseCommunityInviteRequest(c *gin.Context) (*CommunityInviteRequest, error) {
	//POST body params
	var cir CommunityInviteRequest

	if err := c.ShouldBindJSON(&cir); err != nil {
		return nil, err
	}

	return &cir, nil
}
