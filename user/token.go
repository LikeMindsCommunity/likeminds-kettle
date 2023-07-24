package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateTokenRequest struct {
	TokenType string `json:"token_type"`
}

func CreateToken(c *gin.Context) {

	var accessToken string

	// Parsing create token request
	createTokenRequest, _ := parseCreateTokenRequest(c)

	if createTokenRequest != nil && createTokenRequest.TokenType == VTM {
		vtmToken := c.Request.Header.Get(token.HeaderAuthorization)

		if vtmToken == "" {
			utils.TokenAuthError(c, token.ErrorInvalidVTM)
			return
		}

		vtmData, err := token.ExtractVTM(c.Request.Header.Get(token.HeaderAuthorization))

		// If token extraction fails
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Create VTM token
		vtm, tokenErr := token.CreateVTM(vtmData.ApiKey, vtmData.EmailID, vtmData.MobileNo, vtmData.CountryCode)

		// If token creation fails
		if tokenErr != nil {
			utils.GeneralAPIError(c, tokenErr.Error())
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
	utils.GenerateResponse(c, dataResponse)
}

func parseCreateTokenRequest(c *gin.Context) (*CreateTokenRequest, error) {
	//POST body params
	var ctr CreateTokenRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	return &ctr, nil
}
