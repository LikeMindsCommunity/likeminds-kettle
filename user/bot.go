package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/constants"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type BotRequest struct {
	CommunityName string `json:"name"`
}

// CreateBot is used to create bot of a community
func CreateBot(c *gin.Context) {
	Bot(c, utils.POSTMethod)
}

// GetBot is used to get bot of a community
func GetBot(c *gin.Context) {
	Bot(c, utils.GETMethod)
}

// Bot used to create/edit/get bot details of a community
func Bot(c *gin.Context, method int) {

	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	response := GetBotResponse(c, method)
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

// GetBotResponse used to get response when api/user/bot is hit internally
func GetBotResponse(c *gin.Context, method int) *utils.Response {

	//Send request
	var respBytes []byte
	var statusCode int
	var createToken bool
	switch method {
	case utils.GETMethod:

		//Send Request
		respBytes, statusCode = utils.GetRequestResponse(c, utils.CoreService, BotEndpoint, utils.GETRequest, utils.CreateHeaders(c, ""), nil, nil)

	case utils.POSTMethod:

		createToken = true
		botRequest, err := parseBotRequest(c)

		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}

		//Send Request
		respBytes, statusCode = utils.GetRequestResponse(c, utils.CoreService, BotEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, botRequest)

	}

	if respBytes == nil {
		return nil
	}

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return nil
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if createToken {
		userObject := apiCR.Response[ResponseUser].(map[string]interface{})
		userID := userObject[ResponseUserUniqueId].(string)
		userIsGuest := userObject[ResponseUserIsGuest].(bool)

		//Create login and refresh token
		ltm, rtm, err := token.CreateLTMAndRTM(userID, "", token.BETA_AUTH_TOKEN_EXPIRY, token.DEFAULT_TOKEN_EXPIRY, userIsGuest)
		if err != nil {
			//If token creation fails
			utils.GeneralAPIError(c, err.Error())
			return nil
		}
		//Send response with login, refresh token and api/user/login response
		dataResponse[constants.ParamAccessToken] = ltm.AccessToken
		dataResponse[constants.ParamRefreshToken] = rtm.RefreshToken
	}

	return &utils.Response{
		Success: true,
		Data:    dataResponse,
	}
}

func parseBotRequest(c *gin.Context) (*BotRequest, error) {
	//POST body params
	var br BotRequest
	if err := c.ShouldBindBodyWith(&br, binding.JSON); err != nil {
		return nil, err
	}
	return &br, nil
}

func GetUserUniqueIDFromResponse(response *utils.Response) string {
	return response.Data.(map[string]interface{})[ResponseUser].(map[string]interface{})[ResponseUserUniqueId].(string)
}
