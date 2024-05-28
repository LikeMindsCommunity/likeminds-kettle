package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UpdateContentDownloadSetting struct {
	SettingType  string `json:"download_setting_type"`
	SettingTitle string `json:"download_setting_title"`
	IsEnabled    bool   `json:"enabled"`
}

type UpdateContentDownloadSettings struct {
	ContentDownloadSettings []UpdateContentDownloadSetting `json:"content_download_settings"`
}

func GetContentDownloadSettings(c *gin.Context) {
	ContentDownloadSettings(c, utils.GETMethod)
}

func EditContentDownloadSettings(c *gin.Context) {
	ContentDownloadSettings(c, utils.PUTMethod)
}

func ContentDownloadSettings(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request
	switch method {
	case utils.GETMethod:

		// Params to be sent to api/community/fetch_content_download_settings API
		params := map[string]string{
			chatroom.ParamChatroomId: c.Query(chatroom.ParamChatroomId),
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, FetchContentDownloadSettingsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.PUTMethod:

		// Body to be sent in the api/community/update_content_download_settings API internally
		memberRequest, err := parseEditContentDownloadSettingsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, EditContentDownloadSettingsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, memberRequest)

	}
}

func parseEditContentDownloadSettingsRequest(c *gin.Context) (*UpdateContentDownloadSettings, error) {
	//POST body params
	var ucds UpdateContentDownloadSettings

	if err := c.ShouldBindJSON(&ucds); err != nil {
		return nil, err
	}

	return &ucds, nil
}
