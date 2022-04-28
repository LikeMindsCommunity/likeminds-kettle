package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

const VerifyOTPEndPoint = "/api/verify_otp"
const ParamOTP = "otp"
const ResponseProfileExists = "profile_exists"
const ResponseUser = "user"
const ResponseId = "id"

// VerifyOTP is used to verify otp and generate verify token
func VerifyOTP(c *gin.Context) {
	//GET Request params
	otp := c.Query(ParamOTP)
	mobileNo := c.Query(ParamMobileNo)
	countryCode := c.Query(ParamCountryCode)
	if otp == "" || mobileNo == "" || countryCode == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Params to be sent in the api/verify_otp request
	params := map[string]string{
		ParamCountryCode: countryCode,
		ParamMobileNo:    mobileNo,
		ParamOTP:         otp,
	}
	//Create internal API client
	client := api_client.NewAPIClient()
	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + VerifyOTPEndPoint,
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
		//If api/verify_otp success as false
		utils.APIClientError(c, apiCR)
		return
	}
	//If flow succeeds
	profileExists := apiCR.Response[ResponseProfileExists].(bool)
	userID := apiCR.Response[ResponseUser].(map[string]interface{})[ResponseId].(float64)
	//If user exists in our DB, we need to return LTM and RTM
	if profileExists {
		//Create login and refresh token
		ltm, rtm, err := token.CreateLTMAndRTM(utils.FormatFloat(userID, 0))
		if err != nil {
			//If token creation fails
			utils.SomethingWentWrongError(c)
			return
		}
		//Send response with login, refresh token and api/verify_otp response
		dataResponse := apiCR.Response
		dataResponse[token.ParamAccessToken] = ltm.AccessToken
		dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Data:    dataResponse,
		})
		return
	} else {
		//Create verify token
		vtm, err := token.CreateVTM()
		//If token creation fails
		if err != nil {
			utils.SomethingWentWrongError(c)
			return
		}
		//Send response with verify token
		dataResponse := apiCR.Response
		dataResponse[token.ParamAccessToken] = vtm.AccessToken
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Data:    dataResponse,
		})
		return
	}
}
