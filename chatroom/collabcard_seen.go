package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//CollabcardSeen is used to mark a chatroom as seen
func CollabcardSeen(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Params to be sent in the api/collabcard_seen request
	params := map[string]string{
		ParamCollabcardId:   c.Query(ParamCollabcardId),
		ParamCommunityId:    c.Query(ParamCommunityId),
		ParamMemberId:       c.Query(ParamMemberId),
		ParamCollabcardType: c.Query(ParamCollabcardType),
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + CollabcardSeenEndPoint,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
		Params:        params,
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}
