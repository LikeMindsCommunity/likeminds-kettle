package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

func CreateOTMToken(c *gin.Context) {
	// Create verify token
	otm, err := token.CreateOTM(c.GetHeader(utils.HeadersApiKey))

	// If token creation fails
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send response with verify token
	dataResponse := map[string]interface{}{
		token.ParamAccessToken: otm.AccessToken,
	}

	// Generate Response
	utils.GenerateResponse(c, dataResponse)
}
