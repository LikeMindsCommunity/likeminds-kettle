package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type BotRequest struct {
	CommunityName string `json:"community_name"`
}

//CreateBot is used to create bot of a community
func CreateBot(c *gin.Context) {
	Bot(c, utils.POSTMethod)
	return
}

//EditBot is used to edit bot of a community
func EditBot(c *gin.Context) {
	Bot(c, utils.PUTMethod)
	return
}

//GetBot is used to get bot of a community
func GetBot(c *gin.Context) {
	Bot(c, utils.GETMethod)
	return
}

//Bot used when user is signing up and generate login and refresh tokens
func Bot(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Send request
	var respBytes []byte
	var err error
	var createToken bool
	switch method {
	case utils.GETMethod:
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + BotEndpoint,
			CustomHeaders: utils.CreateHeaders(c),
		}
		respBytes, err = client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
	case utils.POSTMethod:
		createToken = true
		br, err := parseBotRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + BotEndpoint,
			Body:          br,
			CustomHeaders: utils.CreateHeaders(c),
		}
		respBytes, err = client.PostRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
	case utils.PUTMethod:
		br, err := parseBotRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + BotEndpoint,
			Body:          br,
			CustomHeaders: utils.CreateHeaders(c),
		}
		respBytes, err = client.PutRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
	}
	if !apiCR.Success {
		//If api/user/login returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	dataResponse := apiCR.Response
	if createToken {
		userID := apiCR.Response[ResponseUser].(map[string]interface{})[ResponseUserUniqueId].(string)
		//Create login and refresh token
		ltm, rtm, err := token.CreateLTMAndRTM(userID)
		if err != nil {
			//If token creation fails
			utils.GeneralAPIError(c, err.Error())
			return
		}
		//Send response with login, refresh token and api/user/login response
		dataResponse[token.ParamAccessToken] = ltm.AccessToken
		dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
	}
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
	return
}

func parseBotRequest(c *gin.Context) (*BotRequest, error) {
	//POST body params
	var br BotRequest
	if err := c.ShouldBindJSON(&br); err != nil {
		return nil, err
	}
	return &br, nil
}
