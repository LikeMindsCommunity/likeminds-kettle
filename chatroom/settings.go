package chatroom

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//GetChatroomSettings is used to fetch the chatroom settings
func GetChatroomSettings(c *gin.Context) {

	//GET Request params
	chatroom_id := c.Query(ParamChatroomId)
	if chatroom_id == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Params to be sent in the api/chatroom/fetch_chatroom_settings request
	params := map[string]string{
		ParamChatroomId: chatroom_id,
	}

	//Create internal API client
	apiClient := api_client.NewAPIClient()
	//Send request

	options := api_client.GetRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + FetchChatroomSettingsEndPoint,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		Params:        params,
	}

	respBytes, err := apiClient.GetRequest(&options)
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
		//If api/chatroom/fetch_settings returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response with api/chatroom/fetch_settings response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}
