package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type FeedNotificationSetting struct {
	NotificationType int  `json:"notification_type"`
	Enabled          bool `json:"enabled"`
}

type EditFeedNotificationSettingsRequest struct {
	NotificationSettings []FeedNotificationSetting `json:"notification_settings"`
}

func GetFeedNotificationSettings(c *gin.Context) {
	CommunityFeedNotificationSettings(c, utils.GETMethod)
}

func EditFeedNotificationSettings(c *gin.Context) {
	CommunityFeedNotificationSettings(c, utils.PUTMethod)
}

func parseEditFeedNotificationSettingsRequest(c *gin.Context) (*EditFeedNotificationSettingsRequest, error) {
	//POST body params
	var efnsr EditFeedNotificationSettingsRequest

	if err := c.ShouldBindJSON(&efnsr); err != nil {
		return nil, err
	}

	return &efnsr, nil
}

func CommunityFeedNotificationSettings(c *gin.Context, method int) {

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
		utils.SendRequest(c, utils.CoreService, FeedNotificationSettingEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	case utils.PUTMethod:

		// Body to be sent in the edit feed notification settings api internally
		editFeedNotificationSettingsRequest, err := parseEditFeedNotificationSettingsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, FeedNotificationSettingEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editFeedNotificationSettingsRequest)
	}

}
