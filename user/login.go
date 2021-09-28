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
	vtm, ok := c.MustGet("vtm").(*token.VerifyTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
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

	//Create login token from the response received in api/user/login
	ltm, err := token.CreateLTM(vtm.VerifiedMobileNo, vtm.CountryCode, userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	token := map[string]string{
		"access_token":  ltm.AccessToken,
		"refresh_token": ltm.RefreshToken,
	}
	c.JSON(http.StatusOK, utils.AuthenticationResponse{
		Success: true,
		Data:    token,
	})
}

func Refresh(c *gin.Context) {
	ltm, ok := c.MustGet("refresh_ltm").(*token.LoginTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
		return
	}

	//Create login token from the response received in api/user/login
	ltm, err := token.CreateLTM(ltm.VerifiedMobileNo, ltm.CountryCode, ltm.UserID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	token := map[string]string{
		"access_token":  ltm.AccessToken,
		"refresh_token": ltm.RefreshToken,
	}
	c.JSON(http.StatusOK, utils.AuthenticationResponse{
		Success: true,
		Data:    token,
	})
}
