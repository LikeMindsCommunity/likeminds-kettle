package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"net/http"
)

type LoginRequest struct {
	UserAcquisitionURL string `json:"user_acquisition_url"`
	LoginJSON          string `json:"login_json"`
	GoogleIDToken      string `json:"google_id_token"`
	LoginType          string `json:"type"`
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	verifyOTPTokenMeta, err := token.ExtractVerifyTokenMeta(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	//TODO - call api/user/login and get response by sending mobile no and country code from verifyOTPTokenMeta
	success := true
	errorMessage := ""
	userID := ""
	if !success {
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	//Create verify token from the response received in api/verify_otp
	tokenDetails, err := token.CreateLoginToken(verifyOTPTokenMeta, userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	token := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}
	c.JSON(http.StatusOK, token)
}
