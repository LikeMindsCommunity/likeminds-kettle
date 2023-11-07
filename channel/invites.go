package channel

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UpdateChannelInviteRequest struct {
	ChannelID    string `json:"channel_id" binding:"required"`
	InviteStatus int    `json:"invite_status"`
}

type UpdateChatroomInviteRequest struct {
	ChatroomID   string `json:"chatroom_id" binding:"required"`
	InviteStatus int    `json:"invite_status"`
}

// Method to fetch channel invites
func GetChannelInvites(c *gin.Context) {
	ChannelInvites(c, utils.GETMethod)
}

// Method to update channel invites
func UpdateChannelInvite(c *gin.Context) {
	ChannelInvites(c, utils.PUTMethod)
}

// ChannelInvites method handles channel invites
func ChannelInvites(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request
	switch method {
	case utils.GETMethod:

		// GET Request params
		channel_type := c.Query(ParamChannelType)

		if channel_type == "1" {
			channel_type = strconv.Itoa(chatroom.NormalChatroomType)

		} else {
			channel_type = strconv.Itoa(chatroom.FeedChatroomType)
		}

		//Params to be sent in the api/chatroom/invites request
		params := map[string]string{
			chatroom.ParamPage:          c.Query(chatroom.ParamPage),
			chatroom.ParamPageSize:      c.Query(chatroom.ParamPageSize),
			chatroom.ParamChatroomTypes: fmt.Sprintf("[%s]", channel_type),
		}

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, ChatroomInvitesEndppoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, false)

	case utils.PUTMethod:

		updateChannelInviteRequest, err := parseUpdateChannelInviteRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		updateChannelInviteRequestBody := UpdateChatroomInviteRequest{
			ChatroomID:   updateChannelInviteRequest.ChannelID,
			InviteStatus: updateChannelInviteRequest.InviteStatus,
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, ChatroomInvitesEndppoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, updateChannelInviteRequestBody)

	}
}

func parseUpdateChannelInviteRequest(c *gin.Context) (*UpdateChannelInviteRequest, error) {
	// PUT body params
	var ucir UpdateChannelInviteRequest

	if err := c.ShouldBindBodyWith(&ucir, binding.JSON); err != nil {
		return nil, err
	}

	return &ucir, nil
}
