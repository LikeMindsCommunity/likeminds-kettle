package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ParticipantRequest struct {
	ChatroomID           interface{}   `json:"chatroom_id"`
	ChatroomParticipants []interface{} `json:"chatroom_participants"`
	UUIDs                []string      `json:"uuids"`
	IsSecret             bool          `json:"is_secret"`
	IsChannelInvite      bool          `json:"is_channel_invite"`
}

type InternalParticipantRequest struct {
	ChatroomID                 interface{}   `json:"chatroom_id"`
	SecretChatroomParticipants []interface{} `json:"secret_chatroom_participants"`
	UUIDs                      []string      `json:"uuids"`
	IsChannelInvite            bool          `json:"is_channel_invite"`
}

type RemoveParticipantRequest struct {
	ChatroomID     interface{}   `json:"chatroom_id"`
	MemberID       interface{}   `json:"member_id,omitempty"`
	UUID           string        `json:"uuid,omitempty"`
	RemovedMembers []interface{} `json:"removed_members"`
	UUIDs          []interface{} `json:"uuids"`
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
	ipr.UUIDs = pr.UUIDs
	ipr.IsChannelInvite = true

	return &ipr
}

func parseParticipantsRequest(c *gin.Context) (*ParticipantRequest, error) {
	// POST body params
	var pr ParticipantRequest

	if err := c.ShouldBindBodyWith(&pr, binding.JSON); err != nil {
		return nil, err
	}

	pr.ChatroomID = utils.ParseInterfaceToString(pr.ChatroomID)

	return &pr, nil
}

func parseRemoveParticipantsRequest(c *gin.Context) (*RemoveParticipantRequest, error) {
	// POST body params
	var rpr RemoveParticipantRequest

	if err := c.ShouldBindBodyWith(&rpr, binding.JSON); err != nil {
		return nil, err
	}

	// parse chatroom_id to string
	rpr.ChatroomID = utils.ParseInterfaceToString(rpr.ChatroomID)

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

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchSecretParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	} else {
		//else, call api/chatroom/fetch_participants_meta api internally

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	}
}

func addParticipantsInternal(c *gin.Context, userId string) {

	participantRequest, err := parseParticipantsRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	is_secret := participantRequest.IsSecret

	if !is_secret {
		//If is_secret is missing or false, call add chatroom participant api internally

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, AddParticipantsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, participantRequest)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, false)

	} else {
		//else, call add secret chatroom participant api internally

		//updated body according to secret participant add request
		addSecretParticipantRequest := updateParticipantsRequest(participantRequest)

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, AddSecretParticipantsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, addSecretParticipantRequest)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, false)
	}
}

func removeParticipantsInternal(c *gin.Context, userId string) {

	removeParticipantRequest, err := parseRemoveParticipantsRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	is_secret := removeParticipantRequest.IsSecret

	if is_secret {

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, RemoveSecretParticipantsEndpoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeParticipantRequest)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, false)

	} else {

		// Updated body according to open participant remove request
		if removeParticipantRequest.MemberID != nil {
			removeParticipantRequest.RemovedMembers = []interface{}{removeParticipantRequest.MemberID}
		}

		if removeParticipantRequest.UUID != "" {
			removeParticipantRequest.UUIDs = []interface{}{removeParticipantRequest.UUID}
		}

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, RemoveOpenParticipantsEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, removeParticipantRequest)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, false)
	}
}
