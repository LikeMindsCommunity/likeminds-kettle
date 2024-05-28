package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ChatroomSettingsRequest struct {
	ID         int    `json:"id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	IsSelected bool   `json:"is_selected" binding:"required"`
}

type EditChatroomSettingsRequest struct {
	ChatroomID       interface{}               `json:"chatroom_id"`
	ChatroomSettings []ChatroomSettingsRequest `json:"chatroom_settings"`
}

// GetChatroomSettings is used to fetch the chatroom settings
func GetChatroomSettings(c *gin.Context) {
	ChatroomSettings(c, utils.GETMethod)
}

// EditChatroomSettings is used to edit the chatroom settings
func EditChatroomSettings(c *gin.Context) {
	ChatroomSettings(c, utils.PUTMethod)
}

// Function used by chatroom setting API's
func ChatroomSettings(c *gin.Context, method int) {

	// Authorize User
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
		// Params to be sent in the fetch chatroom settings api internally
		params := map[string]string{
			ParamChatroomId: c.Query(ParamChatroomId),
		}

		// Params Validation
		if params[ParamChatroomId] == "" {
			//If GET params are missing
			utils.GETQueryParamsMissingError(c)
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, FetchChatroomSettingsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.PUTMethod:

		editChatroomRequest, err := parseEditChatroomSettingsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, EditChatroomSettingsEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editChatroomRequest)
	}

}

func parseEditChatroomSettingsRequest(c *gin.Context) (*EditChatroomSettingsRequest, error) {
	// POST body params
	var ecsr EditChatroomSettingsRequest

	if err := c.ShouldBindBodyWith(&ecsr, binding.JSON); err != nil {
		return nil, err
	}

	if ecsr.ChatroomID != nil {
		ecsr.ChatroomID = utils.ParseInterfaceToString(ecsr.ChatroomID)
	}

	return &ecsr, nil
}
