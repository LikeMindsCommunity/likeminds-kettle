package community

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/channel"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetMemberChannels(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	channel_type := c.Query(channel.ParamChannelType)

	var chatroom_types string

	// Get chatroom types based on channel type
	if channel_type == strconv.Itoa(channel.CHAT_BASED_CHANNEL) {

		chatroom_types = strings.Join([]string{
			strconv.Itoa(chatroom.NormalChatroomType),
			strconv.Itoa(chatroom.AnnouncementChatroomType)},
			",")

	} else if channel_type == strconv.Itoa(channel.FEED_BASED_CHANNEL) {

		chatroom_types = strings.Join([]string{
			strconv.Itoa(chatroom.FeedChatroomType)},
			",")

	} else {

		// If channel type is invalid return error
		utils.GeneralBadRequestError(c, "Invalid channel type sent!")
		return
	}

	requestParams := map[string]string{
		ParamUserId:                 c.Query(ParamUserId),
		chatroom.ParamChatroomTypes: fmt.Sprintf("[%s]", chatroom_types),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, GetMemberChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
