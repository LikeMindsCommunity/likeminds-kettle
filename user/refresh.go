package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

//Refresh to generate new LTM and RTM tokens
func Refresh(c *gin.Context) {
	currentRTM, ok := c.MustGet("rtm").(*token.RefreshTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
		return
	}

	//Create login and refresh token meta from ltm
	ltm, rtm, err := token.CreateLTMAndRTM(currentRTM.VerifiedMobileNo, currentRTM.CountryCode, currentRTM.UserID)
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
}
