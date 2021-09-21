package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"net/http"
)

func VerifyOTP(c *gin.Context) {
	//GET Request params
	otp := c.Query("otp")
	mobileNo := c.Query("mobile_no")
	countryCode := c.Query("country_code")
	if otp == "" || mobileNo == "" || countryCode == "" {
		c.JSON(http.StatusBadRequest, "Params missing")
		return
	}

	//TODO - call api/verify_otp and get response
	success := true
	errorMessage := ""
	profileExists := true

	//Check api/verify_otp success and response
	if !success {
		c.JSON(http.StatusInternalServerError, errorMessage)
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
	c.JSON(http.StatusOK, token)
}
