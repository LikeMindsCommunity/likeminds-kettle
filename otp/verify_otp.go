package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

func VerifyOTP(c *gin.Context) {
	//GET Request params
	otp := c.Query("otp")
	mobileNo := c.Query("mobile_no")
	countryCode := c.Query("country_code")
	if otp == "" || mobileNo == "" || countryCode == "" {
		c.JSON(http.StatusBadRequest, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Query params missing!",
		})
		return
	}

	//TODO - call api/verify_otp and get response
	success := true
	errorMessage := ""
	profileExists := true

	//Check api/verify_otp success and response
	if !success {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: errorMessage,
		})
		return
	}
	if profileExists {

	} else {

	}

	//Create verify token from the response received in api/verify_otp
	tokenDetails, err := token.CreateVerifyOTPToken(mobileNo, countryCode)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	token := map[string]string{
		"access_token": tokenDetails.AccessToken,
	}
	c.JSON(http.StatusOK, utils.AuthenticationResponse{
		Success: true,
		Data:    token,
	})
}
