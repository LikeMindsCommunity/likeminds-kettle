package community

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/channel"
	"github.com/nateshr/likeminds-authentication/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func CommunityMemberChannels(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params
	user_id := c.Param("user_id")
	filterType := []interface{}{}

	switch c.Query(channel.ParamChannelType) {
	case strconv.Itoa(channel.CHAT_BASED_CHANNEL):
		filterType = append(filterType, chatroom.NormalChatroomType, chatroom.AnnouncementChatroomType)

	case strconv.Itoa(channel.FEED_BASED_CHANNEL):
		filterType = append(filterType, chatroom.FeedChatroomType)

	default:
		// If channel type is invalid return error
		utils.GeneralBadRequestError(c, utils.ErrorInvalidChannelType)
		return
	}

	// Params to be sent in the api/fetch_user_chatrooms request
	requestParams := map[string]string{
		ParamState:      c.Query(ParamState),
		ParamPage:       c.Query(ParamPage),
		ParamFilterType: utils.ParseInterfaceListToStringList(filterType),
	}

	// If user_id is digit then send it as user_id else send it as uuid
	if _, err := strconv.Atoi(user_id); err == nil {
		// If user_id is username then fetch user_id from username
		requestParams[ParamUserId] = user_id
	} else {
		requestParams[ParamUUID] = user_id
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, FetchMemberChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
