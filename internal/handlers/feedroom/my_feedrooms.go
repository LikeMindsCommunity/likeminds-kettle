package feedroom

import (
	"strconv"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/chatroom"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// MyFeedrooms is used to fetch all the chatrooms for a user
func MyFeedrooms(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the my chatroom api internally
	params := map[string]string{
		chatroom.ParamPage: c.Query(chatroom.ParamPage),
		chatroom.ParamType: strconv.Itoa(FeedChatroomType),
		chatroom.ParamTag:  c.Query(chatroom.ParamTag),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.MyChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
