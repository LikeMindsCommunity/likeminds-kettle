package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// InitiateSDKEndPoint | togther service user initiate endpoint
const InitiateSDKEndPoint = "/api/sdk/initiate"

// InitiateSDKRequest | user initiate request schema
type InitiateSDKRequest struct {
	UserName     string `json:"user_name"`
	UserUniqueID string `json:"user_unique_id"`
	ImageURL     string `json:"image_url"`
	IsGuest      bool   `json:"is_guest"`
}

//InitiateSDK is used to initiate sdk
func InitiateSDK(c *gin.Context) {
	//POST body bodyParams
	var isr InitiateSDKRequest
	if err := c.ShouldBindJSON(&isr); err != nil {
		//If POST body bodyParams are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}
	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)

	//Hit api/sdk/initiate
	apiClient := api_client.NewAPIClient()
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + InitiateSDKEndPoint,
		Body:          isr,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
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
		//If api/sdk/initiate returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response with login, refresh token and api/sdk/initiate response
	dataResponse := apiCR.Response
	if !ok {
		//If flow succeeds
		userUniqueID := apiCR.Response[user.ResponseUser].(map[string]interface{})[user.ResponseUserUniqueId].(string)
		//Create login and refresh token
		ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID)
		if err != nil {
			//If token creation fails
			utils.GeneralAPIError(c, err.Error())
			return
		}
		dataResponse[token.ParamAccessToken] = ltm.AccessToken
		dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
}
