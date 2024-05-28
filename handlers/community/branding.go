package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func Branding(c *gin.Context) {
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

	getCommunityBrandingEndPoint := fmt.Sprintf(GetCommunityBrandingEndPoint, communityId)

	// Send request
	utils.SendRequest(c, utils.CoreService, getCommunityBrandingEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

}
