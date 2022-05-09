package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateBotRequest struct {
	APIKey string `json:"api_key"`
}

// CreateCommunityBotRequest | create community bot request schema
type CreateCommunityBotRequest struct {
	CommunityName string `json:"community_name" binding:"required"`
}

//CreateBot is used to create bot
func CreateBot(c *gin.Context) {
	//POST body bodyParams
	var isr CreateBotRequest
	if err := c.ShouldBindJSON(&isr); err != nil {
		//If POST body bodyParams are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	apiClient := api_client.NewAPIClient()
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateBotEndpoint,
		Body:          isr,
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
		//If api/user/create_bot returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//If flow succeeds
	userID := apiCR.Response[ResponseUser].(map[string]interface{})[ResponseUserUniqueId].(float64)
	//Create login and refresh token
	ltm, rtm, err := token.CreateLTMAndRTM(utils.FormatFloat(userID, 0))
	if err != nil {
		//If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}
	//Send response with login, refresh token and api/user/create_bot response
	dataResponse := apiCR.Response
	dataResponse[token.ParamAccessToken] = ltm.AccessToken
	dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
	return
}
