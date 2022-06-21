package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

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

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	profileExists := apiCR.Response[ResponseProfileExists].(bool)
	//If user exists in our DB, we need to return LTM and RTM
	if profileExists {
		//Create login and refresh token
		userUniqueID := apiCR.Response[ResponseUser].(map[string]interface{})[ResponseUserUniqueId].(string)
		ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID)
		if err != nil {
			//If token creation fails
			utils.GeneralAPIError(c, err.Error())
			return
		}
		//Send response with login, refresh token and api/verify_otp response
		dataResponse := apiCR.Response
		dataResponse[token.ParamAccessToken] = ltm.AccessToken
		dataResponse[token.ParamRefreshToken] = rtm.RefreshToken

		//Generate Response
		utils.GenerateResponse(c, dataResponse)
	} else {
		//Create verify token
		vtm, err := token.CreateVTM()
		//If token creation fails
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}
		//Send response with verify token
		dataResponse := apiCR.Response
		dataResponse[token.ParamAccessToken] = vtm.AccessToken

		//Generate Response
		utils.GenerateResponse(c, dataResponse)
	}
}
