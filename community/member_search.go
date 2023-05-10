package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func MemberSearch(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Params to be sent in the api/search/member_directory request
	requestParams := map[string]string{
		SearchParam:       c.Query(SearchParam),
		SearchTypeParam:   c.Query(SearchTypeParam),
		ParamPage:         c.Query(ParamPage),
		ParamPageSize:     c.Query(ParamPageSize),
		ParamOrderType:    c.Query(ParamOrderType),
		ParamMemberStates: c.Query(ParamMemberStates),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, UserSearchEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)

}
