package chatroom

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ParticipantRequest struct {
	IsSecret bool `json:"is_secret"`
}

type AddParticipantRequest struct {
	ChatroomID           int64   `json:"chatroom_id"`
	ChatroomParticipants []int64 `json:"chatroom_participants"`
}

type AddSecretParticipantRequest struct {
	ChatroomID                 int64   `json:"chatroom_id"`
	SecretChatroomParticipants []int64 `json:"secret_chatroom_participants"`
}

//AddParticipants is used to add participants in chatroom
func AddParticipants(c *gin.Context) {
	Participants(c, utils.POSTMethod)
}

//GetParticipants is used to get chatroom participants
func GetParticipants(c *gin.Context) {
	Participants(c, utils.GETMethod)
}

//Participatns method handles chatroom participants
func Participants(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Send request
	var respBytes []byte
	var err error
	switch method {
	case utils.GETMethod:

		respBytes = getParticipantsInternal(c, client, response)

	case utils.POSTMethod:

		respBytes = addParticipantsInternal(c, client, response)
	}

	if respBytes == nil {
		return
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
		return
	}
	if !apiCR.Success {
		//If chatroom apis returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}

func parseAddParticipantsRequest(c *gin.Context) (*AddParticipantRequest, error) {
	//POST body params
	var apr AddParticipantRequest

	if err := c.ShouldBindJSON(&apr); err != nil {
		return nil, err
	}

	return &apr, nil
}

func parseAddSecretParticipantsRequest(c *gin.Context) (*AddSecretParticipantRequest, error) {
	//POST body params
	var aspr AddSecretParticipantRequest

	if err := c.ShouldBindJSON(&aspr); err != nil {
		return nil, err
	}

	return &aspr, nil
}

func parseParticipantsRequest(c *gin.Context) (*ParticipantRequest, error) {
	//POST body params
	var pr ParticipantRequest

	if err := c.ShouldBindJSON(&pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func getParticipantsInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	var options api_client.GetRequestOptions

	//GET Request params
	is_secret := c.Query(ParamIsSecret)

	//Params to be sent in the fetch_participants_meta request
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	if is_secret == "" || is_secret == "false" {
		//If is_secret is missing or false, call api/chatroom/fetch_participants_meta api internally

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + FetchParticipantsMetaEndPoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			Params:        params,
		}

	} else {
		//else, call api/chatroom/secret/fetch_participants_meta api internally

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + FetchSecretParticipantsMetaEndPoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			Params:        params,
		}
	}

	respBytes, err := client.GetRequest(&options)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func addParticipantsInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	var options api_client.PostRequestOptions

	participantRequest, err := parseParticipantsRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	is_secret := participantRequest.IsSecret

	if !is_secret {
		//If is_secret is missing or false, call api/chatroom/add api internally

		//Body to be sent in the api/chatroom/add POST request
		addParticipantRequest, err := parseAddParticipantsRequest(c)

		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}

		options = api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + AddParticipantsEndPoint,
			Body:          addParticipantRequest,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}

	} else {
		//else, call api/chatroom/secret/add api internally

		//Body to be sent in the api/chatroom/secret/add POST request
		addSecretParticipantRequest, err := parseAddSecretParticipantsRequest(c)

		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return nil
		}

		options = api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + AddSecretParticipantsEndPoint,
			Body:          addSecretParticipantRequest,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}
	}

	respBytes, err := client.PostRequest(&options)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}
