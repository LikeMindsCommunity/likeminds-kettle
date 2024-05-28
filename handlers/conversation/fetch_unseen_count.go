package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchEventUnseenCount is used to fetch event unseen count
func FetchEventUnseenCount(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchEventUnseenCountEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
