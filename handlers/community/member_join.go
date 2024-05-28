package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type AcceptRejectCommunityJoinRequest struct {
	UUID       string `json:"uuid"`
	IsAccepted bool   `json:"is_accepted"`
}

func AcceptRejectJoinCommunity(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Params to be sent in the api/community_member/join request internally
	acceptRejectJoinCommunity, err := parseAcceptRejectJoinCommunityRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, AcceptRejectCommunityJoinEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, acceptRejectJoinCommunity)
}

func parseAcceptRejectJoinCommunityRequest(c *gin.Context) (*AcceptRejectCommunityJoinRequest, error) {
	//POST body params
	var arcj AcceptRejectCommunityJoinRequest
	if err := c.ShouldBindJSON(&arcj); err != nil {
		return nil, err
	}

	return &arcj, nil
}
