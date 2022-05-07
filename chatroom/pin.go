package chatroom

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PinChatroomRequest struct {
	ChatroomID int64 `json:"chatroom_id"`
	Value      bool  `json:"value"`
	Notify     bool  `json:"notify"`
}

//PinChatroom is used to create a new chatroom
func PinChatroom(c *gin.Context) {

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Create headers from login token
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserID

	//POST body params
	var pcr PinChatroomRequest
	if err := c.ShouldBindJSON(&pcr); err != nil {
		//If POST body params are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	//Create internal API client
	apiClient := api_client.NewAPIClient()

	//Send request
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + PinChatroomEndPoint,
		CustomHeaders: headers,
		Body:          pcr,
	})

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
	}

	if !apiCR.Success {
		//If api/chatroom/pin returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response with api/chatroom/pin response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}
