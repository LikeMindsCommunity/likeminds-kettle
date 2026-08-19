package user

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/token"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func UserSocialLogin(c *gin.Context) {
	// Params to be sent in the api/user/social/login GET request
	params := map[string]string{
		ParamLoginType:        c.Query(ParamLoginType),
		ParamSocialLoginToken: c.Query(ParamSocialLoginToken),
		ParamUserData:         c.Query(ParamUserData),
	}

	// Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UserSocialLoginEndpoint, utils.GETRequest, utils.CreateHeaders(c, ""), params, nil)
	if respBytes == nil {
		return
	}

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	// Send response with login, refresh token and api/user/social/login response
	dataResponse := apiCR.Response

	// Get user info from api response
	userSocialInfo, _ := dataResponse[ResponseUser].(map[string]interface{})

	// Create verified token
	vtm, err := token.CreateVTM(c.GetHeader(utils.HeadersApiKey), userSocialInfo[UserEmail].(string), "", "", c.GetHeader(utils.HeadersPlatformType))

	if err != nil {
		// If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}

	dataResponse[constants.ParamAccessToken] = vtm.AccessToken
	dataResponse[ResponseUser] = map[string]interface{}{
		UserName:     userSocialInfo[UserName],
		UserImageUrl: userSocialInfo[UserImageUrl],
	}

	// Generate response
	utils.GenerateResponse(c, dataResponse, false)
}
