package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetUserFeedMeta(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Access query params and url generation
	userID := c.Param("user_id")

	endpoint := fmt.Sprintf(FetchUserFeedMetaEndPoint, userID)

	utils.SendRequest(c, utils.SwarmService, endpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
