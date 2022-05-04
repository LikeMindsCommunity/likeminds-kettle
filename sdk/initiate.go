package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

const InitiateSDKEndPoint = "/api/sdk/initiate"
const ResponseUser = "user"
const ResponseId = "id"

type InitiateSDKRequest struct {
	UserName     string `json:"user_name"`
	UserUniqueId string `json:"user_unique_id"`
	APIKey       string `json:"api_key"`
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

	apiClient := api_client.NewAPIClient()
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:  apiClient.CoreServiceBaseURL + InitiateSDKEndPoint,
		Body: isr,
		CustomHeaders: utils.CreateHeaders(c),
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

	//If flow succeeds
	userID := apiCR.Response[ResponseUser].(map[string]interface{})[ResponseId].(float64)
	//Create login and refresh token
	ltm, rtm, err := token.CreateLTMAndRTM(utils.FormatFloat(userID, 0))
	if err != nil {
		//If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}
	//Send response with login, refresh token and api/sdk/initiate response
	dataResponse := apiCR.Response
	dataResponse[token.ParamAccessToken] = ltm.AccessToken
	dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
	return
}
