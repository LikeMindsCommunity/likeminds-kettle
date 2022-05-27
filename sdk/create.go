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

// CreateSDKPlatform | create SDK api key platform schema
type CreateSDKPlatform struct {
	Type        int    `json:"type"`
	Package     string `json:"package"`
	Certificate string `json:"certificate"`
}

// CommunityBasicBranding | community basic branding schema
type CommunityBasicBranding struct {
	PrimaryColour string `json:"primary_colour"`
}

// CommunityAdvancedBranding | community advanced branding schema
type CommunityAdvancedBranding struct {
	HeaderColour       string `json:"header_colour"`
	ButtonsIconsColour string `json:"buttons_icons_colour"`
	TextLinksColour    string `json:"text_links_colour"`
}

// CommunityBrandingRequest | create SDK api key platform schema
type CommunityBrandingRequest struct {
	Basic    CommunityBasicBranding    `json:"basic"`
	Advanced CommunityAdvancedBranding `json:"advanced"`
}

// CreateSDKRequest | create SDK api key request schema
type CreateSDKRequest struct {
	CommunityName string                   `json:"name" binding:"required"`
	Branding      CommunityBrandingRequest `json:"branding"`
	Headline      string                   `json:"headline" binding:"required"`
	ImageURL      string                   `json:"image_url"`
	Platform      []CreateSDKPlatform      `json:"platform"`
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
	createCommunityBotRequest := user.BotRequest{CommunityName: createSdkRequest.CommunityName}

	// create internal API client
	apiClient := api_client.NewAPIClient()

	// send internal API request
	createBotResponseBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + user.BotEndpoint,
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
