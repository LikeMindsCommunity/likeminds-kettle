package chatroom

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// FetchShareUrl is used to fetch share url for a specific chatroom
func FetchShareUrl(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch share url api internally
	params := map[string]string{
		ParamChatroomId: c.Query(ParamChatroomId),
		ParamApiType:    strconv.Itoa(SdkApiType),
	}

	// if domain url is present in the request, add it to the params
	domainUrl := c.Query(ParamDomain)
	if domainUrl != "" {
		params[ParamDomain] = domainUrl
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchShareUrlEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
