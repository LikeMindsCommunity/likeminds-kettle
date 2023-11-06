package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Expose method to fetch community configurations for a community
func GetCommunityConfigurations(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	params := map[string]string{
		ParamConfigurationTypes: c.Query(ParamConfigurationTypes),
	}

	//Send Request to api/community/configurations
	utils.SendRequest(c, utils.CoreService, FetchCommunityConfigurationsEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

}
