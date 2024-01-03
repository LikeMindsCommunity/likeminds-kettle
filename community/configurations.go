package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// type communityConfigurationsValue struct {
// }

type UpdateCommunityConfigurationsRequest struct {
	Type  string                 `json:"type" binding:"required"`
	Value map[string]interface{} `json:"value" binding:"required"`
}

func parseCommunityConfigurationsRequest(c *gin.Context) (*UpdateCommunityConfigurationsRequest, error) {

	var uccr UpdateCommunityConfigurationsRequest

	if err := c.ShouldBindJSON(&uccr); err != nil {
		return nil, err
	}

	return &uccr, nil
}

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
	utils.SendRequest(c, utils.CoreService, CommunityConfigurationsEndpoint, utils.GETRequest,
		utils.CreateHeaders(c, userId), params, nil)

}

// Expose method to update community configurations for a community
func UpdateCommunityConfigurations(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Parse request
	uccr, err := parseCommunityConfigurationsRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request to api/community/configurations
	utils.SendRequest(c, utils.CoreService, CommunityConfigurationsEndpoint, utils.PATCHRequest,
		utils.CreateHeaders(c, userId), nil, uccr)

}
