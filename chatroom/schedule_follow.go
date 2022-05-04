package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

const ScheduleFollowEndPoint = "/api/chatroom/schedule_follow"
const ResponseUser = "user"
const ResponseId = "id"

type ScheduleFollowRequest struct {
	ChatroomID         int32 `json:"chatroom_id"`
	ScheduleTime       int64 `json:"schedule_time"`
	ScheduleTimeBefore int64 `json:"schedule_time_before"`
	EndTime            int64 `json:"end_time"`
	EndTimeAfter       int64 `json:"end_time_after"`
}

//ScheduleFollow is used to schedule follow request for particular user
func ScheduleFollow(c *gin.Context) {
	//Check if request has valid login token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//POST body bodyParams
	var isr ScheduleFollowRequest
	if err := c.ShouldBindJSON(&isr); err != nil {
		//If POST body bodyParams are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	apiClient := api_client.NewAPIClient()
	headers := utils.CreateHeaders(c)
	headers[utils.HeadersMemberId] = ltm.UserID
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + ScheduleFollowEndPoint,
		Body:          isr,
		CustomHeaders: headers,
	})
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}
	if !apiCR.Success {
		//If api/chatroom/schedule_follow returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
	return
}
