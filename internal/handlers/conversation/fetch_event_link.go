package conversation

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// FetchEventLink is used to fetch event link
func FetchEventLink(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the fetch event link api internally
	params := map[string]string{
		ParamConversationId: c.Query(ParamConversationId),
	}

	//Params Validation
	if params[ParamConversationId] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EventFetchLinkEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
