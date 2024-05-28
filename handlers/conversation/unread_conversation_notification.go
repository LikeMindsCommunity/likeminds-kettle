package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// UnreadConversationNotification is used to fetch list of unread conversation for notification
func UnreadConversationNotification(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, UnreadConversationNotificationEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
