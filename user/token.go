package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

func CreateToken(c *gin.Context) {

	var accessToken string

	tokenType := c.Query(ParamTokenType)

	if tokenType != "" && tokenType == VTM {
		vtmToken := c.Request.Header.Get(token.HeaderAuthorization)

		if vtmToken == "" {
			utils.TokenAuthError(c, token.ErrorInvalidVTM)
			return
		}

		vtmData, err := token.ExtractVTM(c.Request.Header.Get(token.HeaderAuthorization))

		// If token extraction fails
		if err != nil {
			utils.TokenAuthError(c, err.Error())
			return
		}

		// Create VTM token
		vtm, tokenErr := token.CreateVTM(vtmData.ApiKey, vtmData.EmailID, vtmData.MobileNo, vtmData.CountryCode)

		// If token creation fails
		if tokenErr != nil {
			utils.GeneralBadRequestError(c, tokenErr.Error())
			return
		}

		accessToken = vtm.AccessToken

	} else {
		// Create OTM token
		otm, tokenErr := token.CreateOTM(c.GetHeader(utils.HeadersApiKey))

		// If token creation fails
		if tokenErr != nil {
			utils.GeneralAPIError(c, tokenErr.Error())
			return
		}

		accessToken = otm.AccessToken
	}

	// Send response with verify token
	dataResponse := map[string]interface{}{
		token.ParamAccessToken: accessToken,
	}

	// Generate Response
	utils.GenerateResponse(c, dataResponse, false)
}
