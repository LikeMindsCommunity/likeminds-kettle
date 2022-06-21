package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchShareUrl is used to fetch share url for a specific chatroom
func FetchShareUrl(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Params to be sent in the api/fetch_share_url request
	params := map[string]string{
		ParamChatroomId:  c.Query(ParamChatroomId),
		ParamCommunityId: c.Query(ParamCommunityId),
	}

	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + FetchShareUrlEndPoint,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
		Params:        params,
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
