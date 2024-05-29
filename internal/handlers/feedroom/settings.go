package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// GetFeedroomSettings is used to fetch the feedroom settings
func GetFeedroomSettings(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Params to be sent in the fetch feedroom settings api internally
	params := map[string]string{
		chatroom.ParamChatroomId: c.Query(ParamFeedroomId),
	}

	//Params Validation
	if params[chatroom.ParamChatroomId] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.FetchChatroomSettingsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
