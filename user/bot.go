package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type BotRequest struct {
	CommunityName string `json:"name"`
}

//CreateBot is used to create bot of a community
func CreateBot(c *gin.Context) {
	Bot(c, utils.POSTMethod)
}

//EditBot is used to edit bot of a community
func EditBot(c *gin.Context) {
	Bot(c, utils.PUTMethod)
}

//GetBot is used to get bot of a community
func GetBot(c *gin.Context) {
	Bot(c, utils.GETMethod)
}

//Bot used to create/edit/get bot details of a community
func Bot(c *gin.Context, method int) {
	response := GetBotResponse(c, method)
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

//GetBotResponse used to get response when api/user/bot is hit internally
func GetBotResponse(c *gin.Context, method int) *utils.Response {

	userId := GetRequestingUserId(c)

	if userId == "" {
		return nil
	}

	//Send request
	var respBytes []byte
	var createToken bool
	switch method {
	case utils.GETMethod:

		//Send Request
		respBytes = utils.GetRequestResponse(c, utils.CoreService, BotEndpoint, utils.GETRequest, utils.CreateHeaders(c, ""), nil, nil)

	case utils.POSTMethod:

		createToken = true
		botRequest, err := parseBotRequest(c)

		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}

		//Send Request
		respBytes = utils.GetRequestResponse(c, utils.CoreService, BotEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, botRequest)

	case utils.PUTMethod:

		botRequest, err := parseBotRequest(c)

		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}

		//Send Request
		respBytes = utils.GetRequestResponse(c, utils.CoreService, BotEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, botRequest)

	}

	if respBytes == nil {
		return nil
	}

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes)
	if apiCR == nil {
		return nil
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
			return nil
		}
		//Send response with login, refresh token and api/user/login response
		dataResponse[token.ParamAccessToken] = ltm.AccessToken
		dataResponse[token.ParamRefreshToken] = rtm.RefreshToken
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
