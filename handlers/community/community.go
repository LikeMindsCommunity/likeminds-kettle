package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func Community(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	var communityId string

	community_id, ok := c.Get(ParamCommunityID)

	if !ok {
		sdkAuthResp, err := utils.AuthenticateAPIKeyInternally(utils.CreateHeaders(c, ""), c.GetHeader(utils.HeadersApiKey))

		if err != nil {
			return
		}

		communityId = sdkAuthResp[ParamCommunityID]

	} else {
		communityId = community_id.(string)
	}

	if c.GetHeader(utils.HeadersAcceptVersion) == "v2" {
		params := map[string]string{
			ParamCommunityID: communityId,
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, GetCommunityV2Endpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {
		getCommunityEndPoint := fmt.Sprintf(GetCommunityEndPoint, communityId)

		// Send Request
		utils.SendRequest(c, utils.CoreService, getCommunityEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
	}
}
