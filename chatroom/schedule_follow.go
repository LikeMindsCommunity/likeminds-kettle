package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ScheduleFollowRequest struct {
	ChatroomID int32 `json:"chatroom_id"`
}

//ScheduleFollow is used to schedule follow request for particular user
func ScheduleFollow(c *gin.Context) {
	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}
	//POST body bodyParams
	var sfr ScheduleFollowRequest
	if err := c.ShouldBindJSON(&sfr); err != nil {
		//If POST body bodyParams are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	//Create internal API client
	apiClient := api_client.NewAPIClient()

	//Send request
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + ScheduleFollowEndPoint,
		Body:          sfr,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
	}, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}
