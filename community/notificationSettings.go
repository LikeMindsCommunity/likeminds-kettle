package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type notificationSetting struct {
	NotificationType  int    `json:"notification_type,omitempty"`
	Enabled           *bool  `json:"enabled,omitempty"`
	NotificationState int    `json:"noti_state,omitempty"`
	NotificationTitle string `json:"notification_title,omitempty"`
}

type editNotificationSettingsRequest struct {
	Type                 string                `json:"type"`
	NotificationSettings []notificationSetting `json:"notification_settings"`
}

func parseEditCommunityNotificationSettingsRequest(c *gin.Context) (*editNotificationSettingsRequest, error) {

	var ecnsr editNotificationSettingsRequest

	if err := c.ShouldBindJSON(&ecnsr); err != nil {
		return nil, err
	}

	// Check if notification settings are sent in the request
	if len(ecnsr.NotificationSettings) == 0 {
		return nil, fmt.Errorf(utils.ErrorInvalidNotificationSettings)
	}

	return &ecnsr, nil
}

// Exposed method to fetch Notification Settings for a community
func GetNotificationSettings(c *gin.Context) {
	communityNotificationSettings(c, utils.GETMethod)

}

// Exposed method to edit Notification Settings for a community
func EditNotificationSettings(c *gin.Context) {
	communityNotificationSettings(c, utils.PUTMethod)
}

func communityNotificationSettings(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// If x-accept-version header is not present, then add v1 as default
	if c.GetHeader(utils.HeadersAcceptVersion) == "" {
		c.Request.Header.Add(utils.HeadersAcceptVersion, utils.ApiRevampV1)
	}

	switch method {
	case utils.GETMethod:
		getCommunityNotificationsInternal(c, userId)

	case utils.PUTMethod:
		editCommunityNotificationsInternal(c, userId)

	}
}

func getCommunityNotificationsInternal(c *gin.Context, userId string) {

	// check params
	notificationType := c.Query(ParamType)

	switch notificationType {
	case NotificationTypeChat:

		//Send Request to /api/community/notification_settings
		utils.SendRequest(c, utils.CoreService, ConversationNotificationSettingsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
		return

	case NotificationTypeFeed:

		//Send Request to /api/community/feed_notification_setting
		utils.SendRequest(c, utils.CoreService, FeedNotificationSettingEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
		return

	default:
		utils.GeneralBadRequestError(c, utils.ErrorInvalidNotificationType)
		return
	}

}

func editCommunityNotificationsInternal(c *gin.Context, userId string) {

	// parse and validate request
	ecnsr, err := parseEditCommunityNotificationSettingsRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	switch ecnsr.Type {
	case NotificationTypeChat:

		//Send Request to /api/community/notification_settings
		utils.SendRequest(c, utils.CoreService, ConversationNotificationSettingsEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, ecnsr.NotificationSettings[0])
		return

	case NotificationTypeFeed:

		//Send Request to /api/community/feed_notification_setting
		utils.SendRequest(c, utils.CoreService, FeedNotificationSettingEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, ecnsr)
		return

	default:
		utils.GeneralBadRequestError(c, utils.ErrorInvalidNotificationType)
		return
	}

}
