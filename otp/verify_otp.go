package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

// VerifyOTP is used to verify otp and generate verify token
func VerifyOTP(c *gin.Context) {

	//Params to be sent in the verify request api internally
	params := map[string]string{
		ParamCountryCode: c.Query(ParamCountryCode),
		ParamMobileNo:    c.Query(ParamMobileNo),
		ParamOTP:         c.Query(ParamOTP),
	}

	//Params Validation
	if params[ParamOTP] == "" || params[ParamMobileNo] == "" || params[ParamCountryCode] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	respBytes := utils.GetRequestResponse(c, utils.CoreService, VerifyOTPEndPoint, utils.GETRequest, nil, params, nil)
	if respBytes == nil {
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
		//Send response with login, refresh token and verify otp api response
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
