package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Request Body Struct for Edit conversation notification settings
type EditConversationNotificationSettingsRequest struct {
	NotificationState int `json:"noti_state"`
}

// GetConversationNotificationSettings method to get conversation notification settings
func GetConversationNotificationSettings(c *gin.Context) {
	CommunityNotificationSettings(c, utils.GETMethod)
}

// EditConversationNotificationSettings method to update conversation notification settings
func EditConversationNotificationSettings(c *gin.Context) {
	CommunityNotificationSettings(c, utils.PUTMethod)
}

func parseEditConversationNotificationSettingsRequest(c *gin.Context) (*EditConversationNotificationSettingsRequest, error) {
	//POST body params
	var ecnsr EditConversationNotificationSettingsRequest

	if err := c.ShouldBindJSON(&ecnsr); err != nil {
		return nil, err
	}

	return &ecnsr, nil
}

// CommunityNotificationSettings method to handle community conversation notification settings
func CommunityNotificationSettings(c *gin.Context, method int) {

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
		utils.SendRequest(c, utils.CoreService, ConversationNotificationSettingsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	case utils.PUTMethod:

		// Body to be sent in the edit conversation notification settings api internally
		editConversationNotificationSettingsRequest, err := parseEditConversationNotificationSettingsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ConversationNotificationSettingsEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editConversationNotificationSettingsRequest)
	}

}
