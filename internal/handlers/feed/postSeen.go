package feed

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type PostSeenRequest struct {
	PostIDs []string `json:"post_ids" binding:"required"`
}

// SeenPost is used to mark post as seen for user
func SeenPost(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	// Use Create post body params to create Pending post
	ppsr, err := parsePostSeenRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	if !utils.IsPersonalisedFeedEnabled(utils.GetRedisClientFromContext(c), headers) {
		utils.GeneralBadRequestError(c, utils.PersonalisedFeedDisabledError)
		return
	} else {
		// Send request
		utils.SendRequest(c, utils.SwarmService, PostSeenEndPoint, utils.POSTRequestRawBody, headers, nil, ppsr)
	}

}

func parsePostSeenRequest(c *gin.Context) (*PostSeenRequest, error) {
	//POST body params
	var psr PostSeenRequest
	raw_data, _ := c.GetRawData()

	if err := json.Unmarshal(raw_data, &psr); err != nil {
		return nil, err
	}

	return &psr, nil
}
