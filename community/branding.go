package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utility"
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
		sdkAuthResp, err := utility.AuthenticateAPIKeyInternally(utils.CreateHeaders(c, ""), c.GetHeader(utils.HeadersApiKey))

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
