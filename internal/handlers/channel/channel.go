package channel

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/feedroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// FetchChannel is used to fetch channels
func FetchChannel(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	if c.Query(ParamChannelId) == "" {
		excludedType := "[]"
		filterType := "[]"

		switch c.Query(ParamChannelType) {
		case strconv.Itoa(CHAT_BASED_CHANNEL):
			excludedType = fmt.Sprintf("[%d]", feedroom.FeedChatroomType)

		case strconv.Itoa(FEED_BASED_CHANNEL):
			filterType = fmt.Sprintf("[%d]", feedroom.FeedChatroomType)
		}

		//Params to be sent in the api/chatroom/fetch_all request
		params := map[string]string{
			ParamPage:                  c.Query(ParamPage),
			chatroom.ParamFilterType:   filterType,
			chatroom.ParamExcludedType: excludedType,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.FetchAllChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {
		//else, call api/chatroom/fetch api internally

		//Params to be sent in the api/chatroom/fetch request
		params := map[string]string{
			chatroom.ParamChatroomId: c.Query(ParamChannelId),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.FetchChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}
