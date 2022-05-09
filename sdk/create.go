package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// CreateSDKEndPoint | create SDK togther service endpoint
const CreateSDKEndPoint = "/api/sdk/create"

// CreateSDKBotUserEndPoint | create SDK community bot togther service endpoint
const CreateSDKBotUserEndPoint = "/api/user/create_bot"

// CreateSDKRequestPlatform | create SDK api key platform schema
type CreateSDKRequestPlatform struct {
	Type        int    `json:"type" binding:"required"`
	Package     string `json:"package" binding:"required"`
	Certificate string `json:"certificate" binding:"required"`
}

// CreateSDKRequest | create SDK api key request schema
type CreateSDKRequest struct {
	CommunityName string                     `json:"name" binding:"required"`
	Headline      string                     `json:"headline" binding:"required"`
	BrandColor    string                     `json:"brand_color"`
	ImageURL      string                     `json:"image_url"`
	Platform      []CreateSDKRequestPlatform `json:"platform"`
}

// CreateSDK | returns SDK client API key and community bot's access token
func CreateSDK(request *gin.Context) {
	// post body params
	var createSdkRequest CreateSDKRequest
	if err := request.ShouldBindJSON(&createSdkRequest); err != nil {
		// if post body required params are missing
		utils.POSTBodyParamsMissingError(request)
		return
	}

	// make community bot request
	createCommunityBotRequest := user.CreateCommunityBotRequest{CommunityName: createSdkRequest.CommunityName}

	// create internal API client
	apiClient := api_client.NewAPIClient()

	// send internal API request
	createBotResponseBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateSDKBotUserEndPoint,
		Body:          createCommunityBotRequest,
		CustomHeaders: utils.CreateHeaders(request),
	})
	if err != nil {
		// error in making API request
		utils.GeneralAPIError(request, err.Error())
		return
	}

	// make internal API client
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(createBotResponseBytes, &apiCR)
	if err != nil {
		// error in parsing response data
		utils.GeneralAPIError(request, err.Error())
		return
	}
	if !apiCR.Success {
		// API failed
		request.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	userID := apiCR.Response[ResponseUser].(map[string]interface{})[user.ResponseUserUniqueId].(string)

	// create access tokens for community bot
	ltm, rtm, err := token.CreateLTMAndRTM(userID)
	if err != nil {
		//If token creation fails
		utils.GeneralAPIError(request, err.Error())
		return
	}

	customHeaders := utils.CreateHeaders(request)
	customHeaders[utils.HeadersMemberId] = userID

	// send internal api request
	createSDKResponseBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateSDKEndPoint,
		Body:          createSdkRequest,
		CustomHeaders: customHeaders,
	})
	err = api_client.UnmarshalAPIClientResponse(createSDKResponseBytes, &apiCR)
	if err != nil {
		// error in parsing response data
		utils.GeneralAPIError(request, err.Error())
		return
	}
	if !apiCR.Success {
		// API failed
		request.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	// create API response data
	dataResponse := apiCR.Response
	dataResponse[token.ParamAccessToken] = ltm.AccessToken
	dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
	request.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
	return
}
