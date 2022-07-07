package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//SyncConversation is used to sync conversations
func SyncConversation(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//GET Request params
	is_diff := c.Query(ParamIsDiff)

	if is_diff == "" || is_diff == "false" {
		//If is_diff is missing or false, call api/sync_conversation api internally

		//Params to be sent in the sync conversation api internally
		params := map[string]string{
			ParamChatroomId: c.Query(ParamChatroomId),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, SyncConversationEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	} else {
		//else, call api/sync_conversation_diff api internally

		//Params to be sent in the sync conversation diff api internally
		params := map[string]string{
			ParamIsSynced: c.Query(ParamIsSynced),
			ParamPage:     c.Query(ParamPage),
			ParamPageSize: c.Query(ParamPageSize),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, SyncConversationDiffEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}
