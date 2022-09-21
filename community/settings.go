package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CommunitySetting struct {
	SettingType     string `json:"setting_type"  binding:"required"`
	SettingTitle    string `json:"setting_title"  binding:"required"`
	SettingSubTitle string `json:"setting_sub_title"  binding:"required"`
	IsEnabled       bool   `json:"enabled"  binding:"required"`
}

type EditCommunitySettingsRequest struct {
	CommunitySettings []CommunitySetting `json:"community_settings"`
}

// GetCommunitySettings is used to get community settings
func GetCommunitySettings(c *gin.Context) {
	CommunitySettings(c, utils.GETMethod)
}

// UpdateCommunitySettings is used to get community settings
func UpdateCommunitySettings(c *gin.Context) {
	CommunitySettings(c, utils.PUTMethod)
}

func parseEditCommunitySettingsRequest(c *gin.Context) (*EditCommunitySettingsRequest, error) {
	//POST body params
	var arr EditCommunitySettingsRequest

	if err := c.ShouldBindJSON(&arr); err != nil {
		return nil, err
	}

	return &arr, nil
}

func CommunitySettings(c *gin.Context, method int) {

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
		utils.SendRequest(c, utils.CoreService, FetchCommunitySettings, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	case utils.PUTMethod:

		// Body to be sent in the edit settings api internally
		editCommunitySettingsRequest, err := parseEditCommunitySettingsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, EditCommunitySettings, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editCommunitySettingsRequest)

	}

}
