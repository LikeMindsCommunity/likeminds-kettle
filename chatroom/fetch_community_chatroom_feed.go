package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchCommunityChatroomFeed is used to fetch community chatroom feed
func FetchCommunityChatroomFeed(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//GET Request params
	community_id := c.Query(ParamCommunityId)
	if community_id == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Params to be sent in the api/fetch_community_chatroom_feed request
	params := map[string]string{
		ParamCommunityId: community_id,
		ParamSize:        c.Query(ParamSize),
	}

	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + FetchCommunityChatroomFeedEndPoint,
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
