package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

type LoginRequest struct {
	UserAcquisitionURL string `json:"user_acquisition_url"`
	LoginJSON          string `json:"login_json"`
	GoogleIDToken      string `json:"google_id_token"`
	LoginType          string `json:"type" binding:"required"`
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Invalid JSON Request!",
		})
		return
	}

	verifyOTPTokenMeta, err := token.ExtractVerifyTokenMeta(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Invalid token!",
		})
		return
	}

	//TODO - call api/user/login and get response by sending mobile no and country code from verifyOTPTokenMeta
	success := true
	errorMessage := ""
	userID := "21"
	if !success {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: errorMessage,
		})
		return
	}

	//Create verify token from the response received in api/verify_otp
	tokenDetails, err := token.CreateLoginToken(verifyOTPTokenMeta, userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	token := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}
	c.JSON(http.StatusOK, utils.AuthenticationResponse{
		Success: true,
		Data:    token,
	})
}
