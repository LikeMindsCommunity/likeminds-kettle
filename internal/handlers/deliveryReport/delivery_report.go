package deliveryReport

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// GetDR is used to get delivery report of conversations
func GetDR(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	params := map[string]string{
		ParamChatroomID:      c.Query(ParamChatroomID),
		ParamConversationIDs: c.Query(ParamConversationIDs),
	}

	//Send Request
	utils.SendRequest(c, utils.PandemoniumService, DeliveryReportEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
