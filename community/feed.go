package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GetFeed method will get community feed data
func GetCommunityFeed(c *gin.Context) {
	CommunityFeed(c, utils.GETMethod)
}

// CommunityFeed mthod handles community feed
func CommunityFeed(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	switch method {
	case utils.GETMethod:

		//Params to be sent in the fetch_feed request
		params := map[string]string{
			chatroom.ParamChatroomId:      c.Query(chatroom.ParamChatroomId),
			ParamPinned:                   c.Query(ParamPinned),
			chatroom.ParamScrollDirection: c.Query(chatroom.ParamScrollDirection),
			ParamOrderType:                c.Query(ParamOrderType),
			ParamPage:                     c.Query(ParamPage),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, CommunityFetchFeedEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	}

}
