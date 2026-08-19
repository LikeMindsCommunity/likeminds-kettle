package community

import (
	"strconv"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/channel"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/chatroom"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetMemberChannels(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	channel_type := c.Query(channel.ParamChannelType)

	var chatroom_types []interface{}

	// Get chatroom types based on channel type
	if channel_type == strconv.Itoa(channel.CHAT_BASED_CHANNEL) {

		chatroom_types = append(chatroom_types,
			chatroom.NormalChatroomType,
			chatroom.AnnouncementChatroomType)

	} else if channel_type == strconv.Itoa(channel.FEED_BASED_CHANNEL) {

		chatroom_types = append(chatroom_types,
			chatroom.FeedChatroomType)

	} else {

		// If channel type is invalid return error
		utils.GeneralBadRequestError(c, utils.ErrorInvalidChannelType)
		return
	}

	// Parse Array to String to send in request
	temp_chatroom_types := utils.ParseInterfaceListToStringList(chatroom_types)

	requestParams := map[string]string{
		ParamUserId:                 c.Query(ParamUserId),
		ParamUUID:                   c.Query(ParamUUID),
		chatroom.ParamChatroomTypes: temp_chatroom_types,
		ParamPage:                   c.Query(ParamPage),
		ParamPageSize:               c.Query(ParamPageSize),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, GetMemberChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
