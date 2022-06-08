package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return nil
	}

	//Send request
	var respBytes []byte
	var createToken bool
	switch method {
	case utils.GETMethod:
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + BotEndpoint,
			CustomHeaders: utils.CreateHeaders(c, ""),
		}
		var err error
		respBytes, err = client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return nil
		}
	case utils.POSTMethod:
		createToken = true
		br, err := parseBotRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + BotEndpoint,
			Body:          br,
			CustomHeaders: utils.CreateHeaders(c, ""),
		}
		respBytes, err = client.PostRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return nil
		}
	case utils.PUTMethod:
		br, err := parseBotRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + BotEndpoint,
			Body:          br,
			CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
		}
		respBytes, err = client.PutRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return nil
		}
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
