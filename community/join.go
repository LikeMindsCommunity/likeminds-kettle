package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type QuestionAnswers struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

type CommunityJoinRequest struct {
	QuestionAnswers []QuestionAnswers `json:"question_answers"`
}

func CommunityJoin(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the api/community_member/join POST request
	communityJoinRequest, err := parseCommunityJoinRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CommunityJoinEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, communityJoinRequest)

}

func parseCommunityJoinRequest(c *gin.Context) (*CommunityJoinRequest, error) {

	// POST body params
	var cjr CommunityJoinRequest

	if err := c.ShouldBindJSON(&cjr); err != nil {
		return nil, err
	}

	return &cjr, nil
}
