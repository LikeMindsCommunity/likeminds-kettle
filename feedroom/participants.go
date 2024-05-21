package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddParticipantRequest struct {
	FeedroomID      interface{}   `json:"feedroom_id" binding:"required"`
	Participants    []interface{} `json:"participants"`
	UUIDS           []string      `json:"uuids"`
	IsSecret        bool          `json:"is_secret"`
	IsChannelInvite bool          `json:"is_channel_invite"`
}

type RemoveParticipantRequest struct {
	FeedroomID    interface{} `json:"feedroom_id" binding:"required"`
	ParticipantID interface{} `json:"participant_id"`
	UUID          string      `json:"uuid"`
	IsSecret      bool        `json:"is_secret"`
}

// AddParticipants is used to add participants in feedroom
func AddParticipants(c *gin.Context) {
	Participants(c, utils.POSTMethod)
}

// GetParticipants is used to get feedroom participants
func GetParticipants(c *gin.Context) {
	Participants(c, utils.GETMethod)
}

// GetParticipants is used to get feedroom participants
func RemoveParticipants(c *gin.Context) {
	Participants(c, utils.DELETEMethod)
}

// Participatns method handles feedroom participants
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

func getParticipantsInternal(c *gin.Context, userId string) {

	// GET Request params
	is_secret := c.Query(chatroom.ParamIsSecret)

	//Params to be sent in the fetch_participants_meta request
	params := map[string]string{
		chatroom.ParamChatroomId:      c.Query(ParamFeedroomId),
		chatroom.ParamParticipantName: c.Query(chatroom.ParamParticipantName),
		chatroom.ParamPage:            c.Query(chatroom.ParamPage),
		chatroom.ParamPageSize:        c.Query(chatroom.ParamPageSize),
	}

	if is_secret == "true" {
		//If is_secret is true, call api/chatroom/secret/fetch_participants_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.FetchSecretParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {
		//else, call api/chatroom/fetch_participants_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.FetchParticipantsMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func addParticipantsInternal(c *gin.Context, userId string) {

	addParticipantRequest, err := parseAddParticipantRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	is_secret := addParticipantRequest.IsSecret

	if !is_secret {
		//If is_secret is missing or false, call add chatroom participant api internally

		participantRequest := chatroom.ParticipantRequest{
			ChatroomID:           addParticipantRequest.FeedroomID,
			ChatroomParticipants: addParticipantRequest.Participants,
			UUIDs:                addParticipantRequest.UUIDS,
			IsSecret:             addParticipantRequest.IsSecret,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.AddParticipantsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, participantRequest)
	} else {
		//else, call add secret chatroom participant api internally

		addSecretParticipantRequest := chatroom.InternalParticipantRequest{
			ChatroomID:                 addParticipantRequest.FeedroomID,
			SecretChatroomParticipants: addParticipantRequest.Participants,
			UUIDs:                      addParticipantRequest.UUIDS,
			IsChannelInvite:            true,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.AddSecretParticipantsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, addSecretParticipantRequest)
	}
}

func removeParticipantsInternal(c *gin.Context, userId string) {

	removeParticipantsRequest, err := parseRemoveParticipantRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	is_secret := removeParticipantsRequest.IsSecret

	if is_secret {
		// If is_secret is true, call remove participant from open chatroom api internally

		removeParticipantRequest := chatroom.RemoveParticipantRequest{
			ChatroomID: removeParticipantsRequest.FeedroomID,
			MemberID:   removeParticipantsRequest.ParticipantID,
			UUID:       removeParticipantsRequest.UUID,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.RemoveSecretParticipantsEndpoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeParticipantRequest)

	} else {
		// else call secret participant remove api internally
		removeParticipantRequest := chatroom.RemoveParticipantRequest{
			ChatroomID: removeParticipantsRequest.FeedroomID,
		}

		if removeParticipantsRequest.ParticipantID != nil {
			removeParticipantRequest.RemovedMembers = []interface{}{removeParticipantsRequest.ParticipantID}
		}

		if removeParticipantsRequest.UUID != "" {
			removeParticipantRequest.UUIDs = []interface{}{removeParticipantsRequest.UUID}
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.RemoveOpenParticipantsEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, removeParticipantRequest)
	}
}

func parseAddParticipantRequest(c *gin.Context) (*AddParticipantRequest, error) {
	// POST body params
	var apr AddParticipantRequest

	if err := c.ShouldBindBodyWith(&apr, binding.JSON); err != nil {
		return nil, err
	}

	apr.FeedroomID = utils.ParseInterfaceToString(apr.FeedroomID)

	return &apr, nil
}

func parseRemoveParticipantRequest(c *gin.Context) (*RemoveParticipantRequest, error) {
	// POST body params
	var rpr RemoveParticipantRequest

	if err := c.ShouldBindBodyWith(&rpr, binding.JSON); err != nil {
		return nil, err
	}

	// Parse feedroom_id to string
	rpr.FeedroomID = utils.ParseInterfaceToString(rpr.FeedroomID)

	return &rpr, nil
}
