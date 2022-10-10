package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type EditCommunityDMSettingsRequest struct {
	State            int    `json:"state"`
	Duration         string `json:"duration"`
	NumberInDuration int    `json:"number_in_duration"`
}

func GetCommunityDMSettings(c *gin.Context) {
	CommunityDMSettings(c, utils.GETMethod)
}

func EditCommunityDMSettings(c *gin.Context) {
	CommunityDMSettings(c, utils.PUTMethod)
}

func parseEditCommunityDMSettingsRequest(c *gin.Context) (*EditCommunityDMSettingsRequest, error) {
	//POST body params
	var ecdmsr EditCommunityDMSettingsRequest

	if err := c.ShouldBindJSON(&ecdmsr); err != nil {
		return nil, err
	}

	return &ecdmsr, nil
}

func CommunityDMSettings(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	switch method {
	case utils.GETMethod:

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchCommunityDMSettingsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	case utils.PUTMethod:

		// Body to be sent in the edit community dm settings api internally
		editCommunityDMSettingsRequest, err := parseEditCommunityDMSettingsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, EditCommunityDMSettingsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editCommunityDMSettingsRequest)
	}

}
