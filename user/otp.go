package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/constants"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UserOTPRequest struct {
	OTPMode     string `json:"otp_mode" binding:"required"`
	MobileNo    string `json:"mobile_no"`
	CountryCode int    `json:"country_code"`
	EmailID     string `json:"email_id"`
	IsRetry     bool   `json:"is_retry"`
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
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, UserOTPEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, generateOTPRequest)

	case utils.GETMethod:
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

		// Create verified token
		vtm, err := token.CreateVTM(c.GetHeader(utils.HeadersApiKey), params[ParamEmailID], params[UserMobileNo], params[UserCountryCode])

		if err != nil {
			// If token creation fails
			utils.GeneralAPIError(c, err.Error())
			return
		}

		dataResponse[constants.ParamAccessToken] = vtm.AccessToken

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
