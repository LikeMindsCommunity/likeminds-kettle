package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// Listen Exposed method to open chatroom socket connection
func Listen(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Authorize User
	deviceID := user.GetRequestingUserDeviceId(c)

	// Params to be sent in the api/user/otp GET request
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	// Send request to /chatroom/listen
	utils.SendRequest(c, utils.PandemoniumService, ChatroomListenEndPoint, utils.GETMethod, utils.CreateHeadersFromToken(c, userId, deviceID), params, nil)
}
