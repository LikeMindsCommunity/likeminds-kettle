package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// FetchPreviewUnreadMessagesCount is used to fetch count of unread preview messages
func FetchPreviewUnreadMessagesCount(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch preview unread message count api internally
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchPreviewUnreadMessagesCountEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
