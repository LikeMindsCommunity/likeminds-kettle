package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ParticipantRequest struct {
	ChatroomID           int64         `json:"chatroom_id"`
	ChatroomParticipants []interface{} `json:"chatroom_participants"`
	IsSecret             bool          `json:"is_secret"`
	IsChannelInvite      bool          `json:"is_channel_invite"`
}

type InternalParticipantRequest struct {
	ChatroomID                 int64         `json:"chatroom_id"`
	SecretChatroomParticipants []interface{} `json:"secret_chatroom_participants"`
	IsChannelInvite            bool          `json:"is_channel_invite"`
}

type RemoveParticipantRequest struct {
	ChatroomID     int           `json:"chatroom_id"`
	MemberID       interface{}   `json:"member_id"`
	RemovedMembers []interface{} `json:"removed_members"`
	IsSecret       bool          `json:"is_secret"`
}

// AddParticipants is used to add participants in chatroom
func AddParticipants(c *gin.Context) {
	Participants(c, utils.POSTMethod)
}

// GetParticipants is used to get chatroom participants
func GetParticipants(c *gin.Context) {
	Participants(c, utils.GETMethod)
}

// GetParticipants is used to get chatroom participants
func RemoveParticipants(c *gin.Context) {
	Participants(c, utils.DELETEMethod)
}

// Participatns method handles chatroom participants
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

	case utils.DELETEMethod:

		removeParticipantsInternal(c, userId)
	}
}

func updateParticipantsRequest(pr *ParticipantRequest) *InternalParticipantRequest {
	// POST body params
	var ipr InternalParticipantRequest

	ipr.ChatroomID = pr.ChatroomID
	ipr.SecretChatroomParticipants = pr.ChatroomParticipants
	ipr.IsChannelInvite = pr.IsChannelInvite

	return &ipr
}

func parseParticipantsRequest(c *gin.Context) (*ParticipantRequest, error) {
	// POST body params
	var pr ParticipantRequest

	if err := c.ShouldBindBodyWith(&pr, binding.JSON); err != nil {
		return nil, err
	}

	return &pr, nil
}

func parseRemoveParticipantsRequest(c *gin.Context) (*RemoveParticipantRequest, error) {
	// POST body params
	var rpr RemoveParticipantRequest

	if err := c.ShouldBindBodyWith(&rpr, binding.JSON); err != nil {
		return nil, err
	}

	return &rpr, nil
}

func getParticipantsInternal(c *gin.Context, userId string) {

	// GET Request params
	is_secret := c.Query(ParamIsSecret)

	//Params to be sent in the fetch_participants_meta request
	params := map[string]string{
		ParamChatroomId:      c.Query(ParamChatroomId),
		ParamParticipantName: c.Query(ParamParticipantName),
		ParamPage:            c.Query(ParamPage),
		ParamPageSize:        c.Query(ParamPageSize),
	}

	if is_secret == "true" {
		//If is_secret is true, call api/chatroom/secret/fetch_participants_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchSecretParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {
		//else, call api/chatroom/fetch_participants_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
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

func removeParticipantsInternal(c *gin.Context, userId string) {

	removeParticipantRequest, err := parseRemoveParticipantsRequest(c)

	if err != nil {
		// If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	is_secret := removeParticipantRequest.IsSecret

	if is_secret {
		// If is_secret is true, call remove participant from open chatroom api internally
		utils.SendRequest(c, utils.CoreService, RemoveSecretParticipantsEndpoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeParticipantRequest)

	} else {
		// Updated body according to secret participant remove request
		removeParticipantRequest.RemovedMembers = []interface{}{removeParticipantRequest.MemberID}

		//Send Request
		utils.SendRequest(c, utils.CoreService, RemoveOpenParticipantsEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, removeParticipantRequest)
	}
}
