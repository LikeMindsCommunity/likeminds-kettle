package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//Refresh to generate new LTM and RTM tokens
func Refresh(c *gin.Context) {

	//Check if request has RTM token or not
	currentRTM, ok := c.MustGet(token.ParamRTM).(*token.RefreshTokenMeta)
	if !ok {
		utils.GeneralAPIError(c, utils.ErrorInvalidRTM)
		return
	}

	//Create login and refresh token meta from ltm
	ltm, rtm, err := token.CreateLTMAndRTM(currentRTM.UserUniqueID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}

	//Generate Token Object
	token := map[string]interface{}{
		token.ParamAccessToken:  ltm.AccessToken,
		token.ParamRefreshToken: rtm.RefreshToken,
	}

	//Generate Response
	utils.GenerateResponse(c, token)
}
