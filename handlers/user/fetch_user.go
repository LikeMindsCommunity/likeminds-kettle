package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GetRequestingUserId returns the User Unique ID of user based on request
func GetRequestingUserId(c *gin.Context) string {

	var userUniqueId string = ""

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(constants.ParamLTM).(*constants.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return ""
	}

	userUniqueId = ltm.UserUniqueID

	return userUniqueId
}

func GetBotId(c *gin.Context) string {

	var userUniqueId string = ""

	//Get Platform Type and API Key
	platform_type := c.GetHeader(utils.HeadersPlatformType)
	api_key := c.GetHeader(utils.HeadersApiKey)

	if platform_type == string(utils.PlatformDashboard) && api_key != "" {
		//Call GET api/bot to get bot
		response := GetBotResponse(c, utils.GETMethod)
		if response != nil && response.Success {
			userUniqueId = GetUserUniqueIDFromResponse(response)
		}
	}

	return userUniqueId
}
