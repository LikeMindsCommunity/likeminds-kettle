package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

//Logout is used to blacklist login and refresh tokens and logout user
func Logout(c *gin.Context) {

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Check if request has refresh login token or not
	rtm, ok := c.MustGet(token.ParamRTM).(*token.RefreshTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidRTM)
		return
	}

	//Get redis clients
	client, ok := c.MustGet(cache.ParamRedisClient).(*redis.Client)
	if !ok {
		//If redis client is unavailable
		utils.GeneralAPIError(c, utils.ErrorRedisFailed)
		return
	}

	//Send request internally if DeviceId header exists
	if c.GetHeader(utils.HeadersDeviceId) != "" {

		//Send Request
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, LogoutEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ltm.UserUniqueID), nil, nil)
		if respBytes == nil {
			return
		}

		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return
		}
	}

	//If flow succeeds
	//Blacklist token
	err := cache.BlacklistToken(client, ltm, rtm)
	if err != nil {
		//If token blacklist returns error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Generate Response
	utils.GenerateResponse(c, nil)
}
