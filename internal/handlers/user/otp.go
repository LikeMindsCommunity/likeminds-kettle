package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/handlers/token"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type UserOTPRequest struct {
	OTPMode     string `json:"otp_mode" binding:"required"`
	MobileNo    string `json:"mobile_no"`
	CountryCode int    `json:"country_code"`
	EmailID     string `json:"email_id" binding:"omitempty,email"`
	IsRetry     bool   `json:"is_retry"`
}

type VerifyOTPQuery struct {
	OTPMode     string `form:"otp_mode" binding:"required"`
	MobileNo    string `form:"mobile_no"`
	CountryCode int    `form:"country_code"`
	EmailID     string `form:"email_id" binding:"omitempty,email"`
	OTP         string `form:"otp" binding:"required"`
}

func GenerateUserOTP(c *gin.Context) {

	UserOTP(c, utils.POSTMethod)
}

func VerifyUserOTP(c *gin.Context) {

	UserOTP(c, utils.GETMethod)
}

func UserOTP(c *gin.Context, method int) {

	switch method {
	case utils.POSTMethod:
		// Body to be sent in the user/otp POST method
		generateOTPRequest, err := parseGenerateOTPRequest(c)
		if err != nil {
			// If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, UserOTPEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, generateOTPRequest)

	case utils.GETMethod:

		// Performing query param validation

		if err := parseVerifyOTPQuery(c); err != nil {
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		// Params to be sent in the api/user/otp GET request
		params := map[string]string{
			ParamOTPMode:    c.Query(ParamOTPMode),
			UserMobileNo:    c.Query(UserMobileNo),
			UserCountryCode: c.Query(UserCountryCode),
			ParamEmailID:    c.Query(ParamEmailID),
			ParamOTP:        c.Query(ParamOTP),
		}

		// Send Request
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UserOTPEndpoint, utils.GETRequest, utils.CreateHeaders(c, ""), params, nil)
		if respBytes == nil {
			return
		}

		// Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)

		if apiCR == nil {
			return
		}

		// Send response with login, refresh token and api/user/otp response
		dataResponse := apiCR.Response

		userExists := apiCR.Response[ResponseUserExists].(bool)
		if userExists {
			userObject := apiCR.Response[ResponseUser].(map[string]interface{})
			userUniqueID := userObject[ResponseUserUniqueId].(string)
			userIsGuest := userObject[ResponseUserIsGuest].(bool)

			ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID, "", token.BETA_AUTH_TOKEN_EXPIRY, token.DEFAULT_TOKEN_EXPIRY, userIsGuest, "")
			if err != nil {
				//If token creation fails
				utils.GeneralAPIError(c, err.Error())
				return
			}

			// Set ltm and user_unique_id in context
			ltm.UserUniqueID = userUniqueID
			c.Set(constants.ParamLTM, ltm)

			//Send response with login, refresh token and verify otp api response
			dataResponse[constants.ParamAccessToken] = ltm.AccessToken
			dataResponse[constants.ParamRefreshToken] = rtm.RefreshToken
		} else {
			// Create verified token
			vtm, err := token.CreateVTM(c.GetHeader(utils.HeadersApiKey), params[ParamEmailID], params[UserMobileNo], params[UserCountryCode], c.GetHeader(utils.HeadersPlatformType))

			if err != nil {
				// If token creation fails
				utils.GeneralAPIError(c, err.Error())
				return
			}

			dataResponse[constants.ParamAccessToken] = vtm.AccessToken
		}

		// Generate response
		utils.GenerateResponse(c, dataResponse, false)
	}
}

func parseGenerateOTPRequest(c *gin.Context) (*UserOTPRequest, error) {
	// POST body params
	var uor UserOTPRequest

	if err := c.ShouldBindJSON(&uor); err != nil {
		return nil, err
	}

	return &uor, nil
}

func parseVerifyOTPQuery(c *gin.Context) error {
	var voq VerifyOTPQuery

	if err := c.ShouldBindQuery(&voq); err != nil {
		return err
	}

	return nil
}
