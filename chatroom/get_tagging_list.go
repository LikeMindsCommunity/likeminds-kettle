package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//GetTaggingList is used to fetch the tag members list for a specific chatroom
func GetTaggingList(c *gin.Context) {

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Create headers from login token
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserUniqueID

	//GET Request params
	chatroom_id := c.Query(ParamChatroomId)
	if chatroom_id == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Params to be sent in the api/chatroom/get_tagging_list request
	params := map[string]string{
		ParamChatroomId: chatroom_id,
	}

	//Create internal API client
	apiClient := api_client.NewAPIClient()
	//Send request
	respBytes, err := apiClient.GetRequest(&api_client.GetRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + GetTaggingListEndPoint,
		CustomHeaders: headers,
		Params:        params,
	})

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}
