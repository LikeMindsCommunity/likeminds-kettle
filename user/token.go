package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

func CreateVTMToken(c *gin.Context) {
	// Create verify token
	vtm, err := token.CreateVTM(c.GetHeader(utils.HeadersApiKey))

	// If token creation fails
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send response with verify token
	dataResponse := map[string]interface{}{
		token.ParamAccessToken: vtm.AccessToken,
	}

	// Generate Response
	utils.GenerateResponse(c, dataResponse)
}
