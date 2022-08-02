package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:

		getParticipantsInternal(c, userId)

	case utils.POSTMethod:

		addParticipantsInternal(c, userId)
	}
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

func getParticipantsInternal(c *gin.Context, userId string) {

	//GET Request params
	is_secret := c.Query(ParamIsSecret)

	//Params to be sent in the fetch_participants_meta request
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
	}

	if is_secret == "" || is_secret == "false" {
		//If is_secret is missing or false, call api/chatroom/fetch_participants_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	} else {
		//else, call api/chatroom/secret/fetch_participants_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchSecretParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func addParticipantsInternal(c *gin.Context, userId string) {

	participantRequest, err := parseParticipantsRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	is_secret := participantRequest.IsSecret

	if !is_secret {
		//If is_secret is missing or false, call add chatroom participant api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, AddParticipantsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, participantRequest)
	} else {
		//else, call add secret chatroom participant api internally

		//updated body according to secret participant add request
		addSecretParticipantRequest := updateParticipantsRequest(participantRequest)

		//Send Request
		utils.SendRequest(c, utils.CoreService, AddSecretParticipantsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, addSecretParticipantRequest)
	}
}
