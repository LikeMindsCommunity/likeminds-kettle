package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//CollabcardFollow is used to follow a specific chatroom
func CollabcardFollow(c *gin.Context) {

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
	chatroom_id := c.Query(ParamCollabcardId)
	if chatroom_id == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Params to be sent in the api/collabcard_follow request
	params := map[string]string{
		ParamCollabcardId: chatroom_id,
		ParamMemberId:     c.Query(ParamMemberId),
		ParamValue:        c.Query(ParamValue),
	}

	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + CollabcardFollowEndPoint,
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
