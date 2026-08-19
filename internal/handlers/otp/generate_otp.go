package otp

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// GenerateOTP is used to generate otp
func GenerateOTP(c *gin.Context) {

	//Params to be sent in the generate otp api internally
	params := map[string]string{
		ParamCountryCode: c.Query(ParamCountryCode),
		ParamMobileNo:    c.Query(ParamMobileNo),
	}

	//Params Validation
	if params[ParamCountryCode] == "" || params[ParamMobileNo] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, GenerateOTPEndPoint, utils.GETRequest, nil, params, nil)
}
