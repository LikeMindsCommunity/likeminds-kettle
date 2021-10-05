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

//Login used when user is signing up and generate LTM and RTM tokens
func Login(c *gin.Context) {
	var lr LoginRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Response{
			Success:      false,
			ErrorMessage: "Invalid JSON Request!",
		})
		return
	}
	vtm, ok := c.MustGet("vtm").(*token.VerifyTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.Response{
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
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: errorMessage,
		})
		return
	}

	//Create login and refresh token meta from the response received in api/user/login
	ltm, rtm, err := token.CreateLTMAndRTM(vtm.VerifiedMobileNo, vtm.CountryCode, userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	token := map[string]string{
		"access_token":  ltm.AccessToken,
		"refresh_token": rtm.RefreshToken,
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    token,
	})
}

