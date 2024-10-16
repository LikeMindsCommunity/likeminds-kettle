package deliveryReport

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
	"strings"
)

// GetDR is used to get delivery report of conversations
func GetDR(c *gin.Context) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}
	
	conversationIDs := c.QueryArray(ParamConversationIDs) // Use QueryArray for list of strings
	// Join conversation IDs into a comma-separated string
	conversationIDsStr := strings.Join(conversationIDs, ",")

	params := map[string]string{
		ParamChatroomID:      c.Query(ParamChatroomID),
		ParamConversationIDs: conversationIDsStr,
	}

	//Send Request
	utils.SendRequest(c, utils.PandemoniumService, DeliveryReportEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
