package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//SyncConversation is used to sync conversations
func SyncConversation(c *gin.Context) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	var options api_client.GetRequestOptions

	//GET Request params
	is_diff := c.Query(ParamIsDiff)

	if is_diff == "" || is_diff == "false" {
		//If is_diff is missing or false, call api/sync_conversation api internally

		//Params to be sent in the api/sync_conversation request
		params := map[string]string{
			ParamChatroomId: c.Query(ParamChatroomId),
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + SyncConversationEndPoint,
			CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
			Params:        params,
		}

	} else {
		//else, call api/sync_conversation_diff api internally

		//Params to be sent in the api/sync_conversation_diff request
		params := map[string]string{
			ParamIsSynced: c.Query(ParamIsSynced),
			ParamPage:     c.Query(ParamPage),
			ParamPageSize: c.Query(ParamPageSize),
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + SyncConversationDiffEndPoint,
			CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
			Params:        params,
		}

	}

	respBytes, err := client.GetRequest(&options)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}
