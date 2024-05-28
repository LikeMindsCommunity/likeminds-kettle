package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/token"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/utils"
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
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, VerifyOTPEndPoint, utils.GETRequest, nil, params, nil)
	if respBytes == nil {
		return
	}

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	profileExists := apiCR.Response[ResponseProfileExists].(bool)
	//If user exists in our DB, we need to return LTM and RTM
	if profileExists {
		//Create login and refresh token
		userObject := apiCR.Response[ResponseUser].(map[string]interface{})
		userUniqueID := userObject[ResponseUserUniqueId].(string)
		userIsGuest := userObject[user.ResponseUserIsGuest].(bool)

		ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID, "", token.BETA_AUTH_TOKEN_EXPIRY, token.DEFAULT_TOKEN_EXPIRY, userIsGuest)
		if err != nil {
			//If token creation fails
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Set ltm and user_unique_id in context
		ltm.UserUniqueID = userUniqueID
		c.Set(constants.ParamLTM, ltm)

		//Send response with login, refresh token and verify otp api response
		dataResponse := apiCR.Response
		dataResponse[constants.ParamAccessToken] = ltm.AccessToken
		dataResponse[constants.ParamRefreshToken] = rtm.RefreshToken

		//Generate Response
		utils.GenerateResponse(c, dataResponse, true)
	} else {
		// Create onboarding token
		otm, err := token.CreateOTM(c.GetHeader(utils.HeadersApiKey))

		// If token creation fails
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send response with verify token
		dataResponse := apiCR.Response
		dataResponse[constants.ParamAccessToken] = otm.AccessToken

		// Generate Response
		utils.GenerateResponse(c, dataResponse, false)
	}
}
