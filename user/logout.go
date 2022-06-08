package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/api_client"
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

	headers := utils.CreateHeaders(c, ltm.UserUniqueID)

	if len(headers[utils.HeadersDeviceId].(string)) > 0 {
		//Create internal API client
		apiClient := api_client.NewAPIClient()
		//Send request
		respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
			Url:           apiClient.CoreServiceBaseURL + LogoutEndPoint,
			CustomHeaders: headers,
		})
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes)
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
	//Send response with success as true
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
	})
}
