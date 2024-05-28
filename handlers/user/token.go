package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/token"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/utils"
)

func CreateToken(c *gin.Context) {

	var accessToken string

	tokenType := c.Query(ParamTokenType)

	if tokenType != "" && tokenType == VTM {
		vtmToken := c.Request.Header.Get(constants.HeaderAuthorization)

		if vtmToken == "" {
			utils.TokenAuthError(c, constants.ErrorInvalidVTM)
			return
		}

		vtmData, err := token.ExtractVTM(c.Request.Header.Get(constants.HeaderAuthorization))

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
		constants.ParamAccessToken: accessToken,
	}

	// Generate Response
	utils.GenerateResponse(c, dataResponse, false)
}
