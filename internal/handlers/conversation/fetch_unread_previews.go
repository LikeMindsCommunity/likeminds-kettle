package conversation

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// FetchUnreadPreviews is used to fetch all the unread previews conversation
func FetchUnreadPreviews(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch unread previews api internally
	params := map[string]string{
		ParamChatroomId:  c.Query(ParamChatroomId),
		ParamPaginatedBy: c.Query(ParamPaginatedBy),
		ParamPage:        c.Query(ParamPage),
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchUnreadPreviewsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}
