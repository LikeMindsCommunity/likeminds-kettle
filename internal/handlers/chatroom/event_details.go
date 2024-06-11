package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// Exposed method to fetch event meta
func FetchEventMeta(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	params := map[string]string{
		ParamPastEvents: c.Query(ParamPastEvents),
	}

	// Send request to /api/chatroom/event/fetch_all_meta
	utils.SendRequest(c, utils.CoreService, FetchEventMetaEndPoint, utils.GETMethod, utils.CreateHeaders(c, userId), params, nil)
}

// Exposed method to fetch event links
func FetchEventLinks(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	params := map[string]string{
		ParamApiType:    c.Query(ParamApiType),
		ParamChatroomId: c.Query(ParamChatroomId),
		ParamIsEditMode: c.Query(ParamIsEditMode),
	}

	// Send request to /api/chatroom/event/fetch_link
	utils.SendRequest(c, utils.CoreService, FetchEventLinksEndPoint, utils.GETMethod, utils.CreateHeaders(c, userId), params, nil)
}

// Exposed method to fetch unseen event count
func FetchEventUnseenCount(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request to /api/chatroom/event/fetch_unseen_count
	utils.SendRequest(c, utils.CoreService, FetchEventUnseenCountEndPoint, utils.GETMethod, utils.CreateHeaders(c, userId), nil, nil)
}
