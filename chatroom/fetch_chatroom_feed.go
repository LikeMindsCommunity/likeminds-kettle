package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchChatroomFeed is used to fetch a specific chatroom feed
func FetchChatroomFeed(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Params to be sent in the /api/v1/fetch_chatroom_feed request
	params := map[string]string{
		ParamCommunityId:     c.Query(ParamCommunityId),
		ParamActive:          c.Query(ParamActive),
		ParamScrollDirection: c.Query(ParamScrollDirection),
		ParamChatroomId:      c.Query(ParamChatroomId),
	}

	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + FetchChatroomFeedEndPoint,
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
