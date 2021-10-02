package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

// VerifyOTP is used to verify otp and generate VTM Token
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
	userId := "21"

	//Check api/verify_otp success and response
	if !success {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: errorMessage,
		})
		return
	}
	if profileExists {
		//Create login and refresh token meta from the response received in api/verify_otp
		ltm, rtm, err := token.CreateLTMAndRTM(mobileNo, countryCode, userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}
		token := map[string]string{
			"access_token":  ltm.AccessToken,
			"refresh_token": rtm.RefreshToken,
		}
		c.JSON(http.StatusOK, utils.AuthenticationResponse{
			Success: true,
			Data:    token,
		})
	} else {
		//Create verify token meta from the response received in api/verify_otp
		vtm, err := token.CreateVTM(mobileNo, countryCode)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		token := map[string]string{
			"access_token": vtm.AccessToken,
		}
		c.JSON(http.StatusOK, utils.AuthenticationResponse{
			Success: true,
			Data:    token,
		})
	}
}
