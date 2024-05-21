package channel

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ChannelSettingsRequest struct {
	SettingType string `json:"setting_type" binding:"required"`
	Enabled     bool   `json:"enabled" binding:"required"`
}

type UpdateUserChannelSettingsRequest struct {
	ChannelSettings []ChannelSettingsRequest `json:"channel_settings" binding:"required"`
}

// Method to fetch channel invites
func GetUserChannelSettings(c *gin.Context) {
	ChannelUserSettings(c, utils.GETMethod)
}

// Method to update channel invites
func UpdateUserChannelSettings(c *gin.Context) {
	ChannelUserSettings(c, utils.PUTMethod)
}

func ChannelUserSettings(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Get URL params
	channel_id := c.Param(ParamChannelId)
	participant_uuid := c.Param(ParamParticipantUUID)

	UserChannelEndpoint := fmt.Sprintf(UserChannelSettingsEndpoint, channel_id, participant_uuid)

	// Send request
	switch method {
	case utils.GETMethod:

		params := map[string]string{
			"setting_types": c.Query("setting_types"),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, UserChannelEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.PUTMethod:

		updateUserChannelSettingsRequest, err := parseUpdateUserChannelSettingsReqest(c)
		if err != nil {
			// If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, UserChannelEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, updateUserChannelSettingsRequest)

	}
}

func parseUpdateUserChannelSettingsReqest(c *gin.Context) (*UpdateUserChannelSettingsRequest, error) {
	var uucsr UpdateUserChannelSettingsRequest

	if err := c.ShouldBindJSON(&uucsr); err != nil {
		return nil, err
	}

	return &uucsr, nil
}
