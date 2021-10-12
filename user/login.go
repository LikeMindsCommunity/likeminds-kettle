package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

type LoginRequest struct {
	UserAcquisitionURL string `json:"user_acquisition_url"`
	LoginJSON          string `json:"login_json"`
	GoogleIDToken      string `json:"google_id_token"`
	LoginType          string `json:"type" binding:"required"`
}

//Login used when user is signing up and generate login and refresh tokens
func Login(c *gin.Context) {
	//Check if request has valid verify token or not
	vtm, ok := c.MustGet("vtm").(*token.VerifyTokenMeta)
	if !ok {
		//If token is not available
		utils.SomethingWentWrongError(c)
		return
	}
	//POST body params
	var lr LoginRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		//If POST body params are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	//Params to be sent in the api/user/login
	params, err := utils.RequestParamsToMap(lr)
	if err != nil {
		//If mapping fails
		utils.SomethingWentWrongError(c)
	}
	//Create internal API client
	client := api_client.NewAPIClient()
	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + "/api/user/login",
		Params:        params,
		CustomHeaders: nil,
	}
	//Send request
	respBytes, err := client.GetRequest(&options)
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
		utils.SomethingWentWrongError(c)
	}

	if !apiCR.Success {
		//If api/user/login returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	emailExists := apiCR.Response["email_exists"].(bool)
	mobileNo := vtm.VerifiedMobileNo
	countryCode := vtm.CountryCode
	if emailExists {
		//Merge account case
		//Create verify tokenResponse meta from the response received in VTM
		vtm, err := token.CreateVTM(mobileNo, countryCode)
		if err != nil {
			//If token creation fails
			utils.SomethingWentWrongError(c)
			return
		}
		//Send response with verify token and api/user/login response
		dataResponse := apiCR.Response
		dataResponse["access_token"] = vtm.AccessToken
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Data:    dataResponse,
		})
		return
	} else {
		//New user login case
		//Get user ID from api/user/login response
		userID := apiCR.Response["user"].(map[string]interface{})["id"].(float64)
		//Create login and refresh token
		ltm, rtm, err := token.CreateLTMAndRTM(mobileNo, countryCode, userID)
		if err != nil {
			//If token creation fails
			utils.SomethingWentWrongError(c)
			return
		}
		//Send response with login, refresh token and api/user/login response
		dataResponse := apiCR.Response
		dataResponse["access_token"] = ltm.AccessToken
		dataResponse["refresh_token"] = rtm.RefreshToken
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Data:    dataResponse,
		})
		return
	}
}
