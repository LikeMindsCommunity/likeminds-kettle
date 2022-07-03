package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//GetRequestingUserId returns the User Unique ID of user based on request
func GetRequestingUserId(c *gin.Context) string {

	var userUniqueId string = ""

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return ""
	}

	userUniqueId = ltm.UserUniqueID

	platform_type := c.GetHeader(utils.HeadersPlatformType)

	if platform_type == string(utils.PlatformDashboard) {
		//Call GET api/bot to get bot
		response := GetBotResponse(c, utils.GETMethod)
		if response != nil {
			userUniqueId = GetUserUniqueIDFromResponse(response)
		}
	}

	return userUniqueId
}
