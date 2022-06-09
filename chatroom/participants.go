package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ParticipantRequest struct {
	ChatroomID           int64   `json:"chatroom_id"`
	ChatroomParticipants []int64 `json:"chatroom_participants"`
	IsSecret             bool    `json:"is_secret"`
}

type InternalParticipantRequest struct {
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
	utils.ParseResponse(c, respBytes)
}

func updateParticipantsRequest(pr *ParticipantRequest) *InternalParticipantRequest {
	//POST body params
	var ipr InternalParticipantRequest

	ipr.ChatroomID = pr.ChatroomID
	ipr.SecretChatroomParticipants = pr.ChatroomParticipants

	return &ipr
}

func parseParticipantsRequest(c *gin.Context) (*ParticipantRequest, error) {
	//POST body params
	var pr ParticipantRequest

	if err := c.ShouldBindBodyWith(&pr, binding.JSON); err != nil {
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

		options = api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + AddParticipantsEndPoint,
			Body:          participantRequest,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}

	} else {
		//else, call api/chatroom/secret/add api internally

		//updated body according to secret participant add request
		addSecretParticipantRequest := updateParticipantsRequest(participantRequest)

		options = api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + AddSecretParticipantsEndPoint,
			Body:          addSecretParticipantRequest,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}
