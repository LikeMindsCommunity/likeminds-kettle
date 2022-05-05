package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

const CreateSDKEndPoint = "/api/sdk/create"
const CreateSDKBotUserEndPoint = "/api/user/create_bot"

type CreateSDKRequestPlatform struct {
	Type        int    `json:"type"`
	Package     string `json:"package"`
	Certificate string `json:"certificate"`
}

type CreateSDKRequest struct {
	CommiunityName string                     `json:"community_name"`
	Headline       string                     `json:"headline"`
	BrandColour    string                     `json:"brand_colour"`
	ImageURL       string                     `json:"image_url"`
	Platform       []CreateSDKRequestPlatform `json:"platform"`
}

func CreateSDK(request *gin.Context) {
	bearerToken := request.GetHeader("Authorization")
	token, err := token.VerifyToken(bearerToken)
	if err != nil || token == nil {
		utils.TokenAuthError(request, err.Error())
		return
	}

	var createSdkRequest CreateSDKRequest
	if err := request.ShouldBindJSON(&createSdkRequest); err != nil {
		utils.POSTBodyParamsMissingError(request)
		return
	}

	apiClient := api_client.NewAPIClient()
	createBotResponseBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateSDKBotUserEndPoint,
		Body:          request.Request.Body,
		CustomHeaders: utils.CreateHeaders(request),
	})
	if err != nil {
		utils.GeneralAPIError(request, err.Error())
		return
	}

	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(createBotResponseBytes, &apiCR)
	if err != nil {
		utils.GeneralAPIError(request, err.Error())
		return
	}
	if !apiCR.Success {
		request.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	userID := apiCR.Response[ResponseUser].(map[string]interface{})[user.ResponseUserUniqueId].(float64)

	customHeaders := utils.CreateHeaders(request)
	customHeaders[utils.HeadersMemberId] = userID

	createSDKResponseBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateSDKEndPoint,
		Body:          request,
		CustomHeaders: customHeaders,
	})
	err = api_client.UnmarshalAPIClientResponse(createSDKResponseBytes, &apiCR)
	if err != nil {
		utils.GeneralAPIError(request, err.Error())
		return
	}
	if !apiCR.Success {
		request.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	dataResponse := apiCR.Response
	request.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
	return
}
