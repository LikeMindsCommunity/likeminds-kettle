package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type ChatroomFollowRequest struct {
	CollabcardId interface{} `json:"collabcard_id"`
	ChatroomId   interface{} `json:"chatroom_id"`
	MemberId     interface{} `json:"member_id"`
	UUID         string      `json:"uuid"`
	Value        *bool       `json:"value"`
}

// ChatroomFollow is used to follow a specific chatroom
func ChatroomFollow(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the api/collabcard_follow request
	params := map[string]string{
		ParamCollabcardId: c.Query(ParamCollabcardId),
		ParamMemberId:     c.Query(ParamMemberId),
		ParamUUID:         c.Query(ParamUUID),
		ParamValue:        c.Query(ParamValue),
	}

	chatroomFollowRequest, err := parseChatroomFollowRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// If collabcard_id is present in query params, then use it and skip body params
	if params[ParamCollabcardId] == "" && chatroomFollowRequest != nil {
		params[ParamCollabcardId] = utils.ParseInterfaceToString(chatroomFollowRequest.CollabcardId)
		params[ParamMemberId] = utils.ParseInterfaceToString(chatroomFollowRequest.MemberId)
		params[ParamUUID] = chatroomFollowRequest.UUID
		params[ParamValue] = utils.ParseInterfaceToString(chatroomFollowRequest.Value)
	}

	// If collabcard_id or chatroom_id is missing
	if params[ParamCollabcardId] == "" {

		utils.GeneralBadRequestError(c, "collabcard_id or chatroom_id is missing")

		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CollabcardFollowEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}

func parseChatroomFollowRequest(c *gin.Context) (*ChatroomFollowRequest, error) {

	// If body params are missing then return nil
	if c.Request.Body == nil || c.Request.ContentLength == 0 {
		return nil, nil
	}

	//POST body params
	var cfr ChatroomFollowRequest

	if err := c.ShouldBindJSON(&cfr); err != nil {
		return nil, err
	}

	// If chatroom_id is present, then pass it in collabcard_id
	if cfr.ChatroomId != nil {
		cfr.CollabcardId = cfr.ChatroomId
	}

	return &cfr, nil
}
