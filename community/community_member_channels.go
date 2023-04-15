package community

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/channel"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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
		ParamUserId:     user_id,
		ParamState:      c.Query(ParamState),
		ParamPage:       c.Query(ParamPage),
		ParamFilterType: utils.ParseArrayToString(filterType),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, FetchMemberChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
}
