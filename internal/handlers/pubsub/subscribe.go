package pubsub

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// Subscribe Exposed method to open socket connection and subscribe to topic
func Subscribe(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Authorize User
	deviceID := user.GetRequestingUserDeviceId(c)

	// Params to be sent in the api/user/otp GET request
	params := map[string]string{
		ParamTopic: c.Query(ParamTopic),
	}

	// Send request to /chatroom/listen
	utils.SendRequest(c, utils.PandemoniumService, SubscribeEndPoint, utils.GETMethod, utils.CreateHeadersFromToken(c, userId, deviceID), params, nil)
}
