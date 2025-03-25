package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type RemoveCohortMemberRequest struct {
	CohortID int    `json:"cohort_id" binding:"required"`
	MemberID int    `json:"user_id"`
	UUID     string `json:"uuid"`
}

func RemoveCohortMember(c *gin.Context) {
	CohortMember(c, utils.DELETEMethod)
}

func CohortMember(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Send request
	switch method {
	case utils.DELETEMethod:

		// Body to be sent in the remove cohort member internally
		removeCohortMemberRequest, err := parseRemoveCohortMemberRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, RemoveCohortMemberEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, removeCohortMemberRequest)

		// delete cached user feed access rights
		user.DeleteAccessDataAgainstUserIdAndAccessTypeFromCache(utils.GetRedisClientFromContext(c), removeCohortMemberRequest.UUID)

	}
}

func parseRemoveCohortMemberRequest(c *gin.Context) (*RemoveCohortMemberRequest, error) {

	// POST body params
	var rcm RemoveCohortMemberRequest

	if err := c.ShouldBindJSON(&rcm); err != nil {
		return nil, err
	}

	return &rcm, nil
}
