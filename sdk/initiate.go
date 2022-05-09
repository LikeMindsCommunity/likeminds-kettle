package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

const InitiateSDKEndPoint = "/api/sdk/initiate"
const ResponseUser = "user"
const ResponseId = "id"

type InitiateSDKRequest struct {
	UserName     string `json:"user_name" binding:"required"`
	UserUniqueId string `json:"user_unique_id"`
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

	headers := utils.CreateHeaders(c)
	headers[utils.HeadersApiKey] = c.GetHeader(utils.HeadersApiKey)

	apiClient := api_client.NewAPIClient()
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + InitiateSDKEndPoint,
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
		//If api/sdk/initiate returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//If flow succeeds
	userID := apiCR.Response[ResponseUser].(map[string]interface{})[user.ResponseUserUniqueId].(string)
	//Create login and refresh token
	ltm, rtm, err := token.CreateLTMAndRTM(userID)
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
