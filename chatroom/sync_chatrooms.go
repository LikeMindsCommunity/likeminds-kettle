package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//SyncChatrooms is used to fetch data for chatroom syncing
func SyncChatrooms(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//GET Request params
	is_diff := c.Query(ParamIsDiff)

	if is_diff == "" || is_diff == "false" {
		//If is_diff is missing or false, call api/sync_chatrooms api internally

		//Params to be sent in the sync chatroom api internally
		params := map[string]string{
			ParamPage:           c.Query(ParamPage),
			ParamPageSize:       c.Query(ParamPageSize),
			ParamChatroomStatus: c.Query(ParamChatroomStatus),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, SyncChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	} else {
		//else, call api/sync_chatrooms_diff api internally

		//Params to be sent in the sync chatroom diff api internally
		params := map[string]string{
			ParamPage:     c.Query(ParamPage),
			ParamPageSize: c.Query(ParamPageSize),
			ParamIsSynced: c.Query(ParamIsSynced),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, SyncChatroomsDiffEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}
