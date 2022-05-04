package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

//Logout is used to blacklist login and refresh tokens and logout user
func Logout(c *gin.Context) {
	//Check if request has valid login token or not
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

	//Create headers from login token
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserID
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
	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
	}

	if !apiCR.Success {
		//If api/user/logout returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	//Blacklist token
	err = cache.BlacklistToken(client, ltm, rtm)
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
