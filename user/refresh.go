package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

//Refresh to generate new LTM and RTM tokens
func Refresh(c *gin.Context) {
	currentRTM, ok := c.MustGet(token.ParamRTM).(*token.RefreshTokenMeta)
	if !ok {
		utils.SomethingWentWrongError(c)
		return
	}

	//Create login and refresh token meta from ltm
	ltm, rtm, err := token.CreateLTMAndRTM(currentRTM.UserID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	token := map[string]string{
		token.ParamAccessToken:  ltm.AccessToken,
		token.ParamRefreshToken: rtm.RefreshToken,
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    token,
	})
}
